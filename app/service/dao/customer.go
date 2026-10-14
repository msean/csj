package dao

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"

	"gorm.io/gorm"
)

type customerDao struct{}

func newCustomerDao() *customerDao {
	return &customerDao{}
}

func (dao *customerDao) Update(db *gorm.DB, object model.Customer) (err error) {
	// 用事务保证：查旧 Credit → 计算 delta → 同步 Credit 与 Debt 的原子性
	return db.Transaction(func(tx *gorm.DB) error {
		// 查询旧的 Credit，用于计算 Credit 变化差额
		var oldCredit float64
		if err = tx.Model(&model.Customer{}).Where("uid = ?", object.UID).
			Limit(1).Pluck("credit", &oldCredit).Error; err != nil {
			return err
		}

		// 前端直接传入 Credit，Debt 同步等量调整
		delta := object.Credit - oldCredit
		return tx.Model(&model.Customer{}).
			Where("uid = ?", object.UID).
			Updates(map[string]interface{}{
				"name":   object.Name,
				"phone":  object.Phone,
				"remark": object.Remark,
				"addr":   object.Addr,
				"credit": object.Credit,
				"debt":   gorm.Expr("debt + ?", delta),
			}).Error
	})
}

func (dao *customerDao) NewTempCustomer(db *gorm.DB, ownerUser string) (err error) {
	c := model.Customer{
		Name:      common.CashCustomerName,
		Remark:    common.CashCustomerName,
		OwnerUser: ownerUser,
	}
	return utils.CreateObj(db, &c)
}

func (dao *customerDao) List(db *gorm.DB, ownerUser string, condtions request.CustomerListReq) (customers []model.Customer, err error) {
	conds := []utils.Cond{}
	conds = append(conds, []utils.Cond{
		utils.NewWhereCond("owner_user", ownerUser),
		utils.NewOrderCond("credit DESC"),
	}...)
	if condtions.SearchKey != "" {
		conds = append(conds, utils.NewOrLikeCond(condtions.SearchKey, utils.LikeTypeBetween, "name", "phone"))
	}
	err = utils.Find(db, &customers, conds...)
	return
}

func (dao *customerDao) FromUUID(db *gorm.DB, uuid string) (customers model.Customer, err error) {
	err = utils.First(db, &customers, utils.WhereUIDCond(uuid))
	return
}

// IncrOrderStats 新建有效订单时增加客户统计
func (dao *customerDao) IncrOrderStats(db *gorm.DB, userUUID string, totalAmount float64, creditAmount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Updates(map[string]interface{}{
			"order_count":  gorm.Expr("order_count + 1"),
			"total_amount": gorm.Expr("total_amount + ?", totalAmount),
			"credit":       gorm.Expr("credit + ?", creditAmount),
		}).Error
	return
}

// DecrOrderStats 作废/退款/退货时减少客户统计
func (dao *customerDao) DecrOrderStats(db *gorm.DB, userUUID string, totalAmount float64, creditAmount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Updates(map[string]interface{}{
			"order_count":  gorm.Expr("GREATEST(order_count - 1, 0)"),
			"total_amount": gorm.Expr("total_amount - ?", totalAmount),
			"credit":       gorm.Expr("credit - ?", creditAmount),
		}).Error
	return
}

// UpdateOrderStatsDiff 订单更新时调整客户统计差额
func (dao *customerDao) UpdateOrderStatsDiff(db *gorm.DB, userUUID string, deltaTotal float64, deltaDebt float64) (err error) {
	updates := map[string]interface{}{
		"total_amount": gorm.Expr("total_amount + ?", deltaTotal),
	}
	if deltaDebt != 0 {
		updates["credit"] = gorm.Expr("credit + ?", deltaDebt)
	}
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Updates(updates).Error
	return
}

// DecrDebt 还款时扣减客户欠款
func (dao *customerDao) DecrDebt(db *gorm.DB, userUUID string, amount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Update("credit", gorm.Expr("credit - ?", amount)).Error
	return
}

// IncrDebt 撤销还款时恢复客户欠款
func (dao *customerDao) IncrDebt(db *gorm.DB, userUUID string, amount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Update("credit", gorm.Expr("credit + ?", amount)).Error
	return
}

func (dao *customerDao) Pay(db *gorm.DB, userUUID string, amount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Update("paid_debt", gorm.Expr("COALESCE(paid_debt, 0) + ?", amount)).Error
	return
}

// DecrPaidDebt 撤销原始债务还款时扣减 PaidDebt
func (dao *customerDao) DecrPaidDebt(db *gorm.DB, userUUID string, amount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Update("paid_debt", gorm.Expr("GREATEST(COALESCE(paid_debt, 0) - ?, 0)", amount)).Error
	return
}

// GetDebtPaidDebt 获取客户当前 Debt 与 PaidDebt，用于判断原始债务还款进度
// func (dao *customerDao) GetDebtPaidDebt(db *gorm.DB, userUUID string) (debt, paidDebt float64, err error) {
// 	var c model.Customer
// 	if err = db.Model(&model.Customer{}).Select("debt, paid_debt").
// 		Where("uid = ?", userUUID).First(&c).Error; err != nil {
// 		return
// 	}
// 	return c.Debt, c.PaidDebt, nil
// }

// OwnerStatsResult 商户下所有客户的累计统计
type OwnerStatsResult struct {
	TotalAmount float64 `gorm:"column:total_amount"`
	TotalDebt   float64 `gorm:"column:total_debt"`
}

// SumOwnerStats 统计商户下所有客户的累计销售额（total_amount 之和）和总赊欠额（debt 之和）
func (dao *customerDao) SumOwnerStats(db *gorm.DB, ownerUser string) (stats OwnerStatsResult, err error) {
	err = db.Model(&model.Customer{}).
		Select("COALESCE(SUM(total_amount), 0) AS total_amount, COALESCE(SUM(credit), 0) AS total_debt").
		Where("owner_user = ?", ownerUser).
		Scan(&stats).Error
	return
}

// DebtCustomer 有赊欠的客户（只取需要的字段）
type DebtCustomer struct {
	UID  string  `gorm:"column:uid" json:"uuid"`
	Name string  `gorm:"column:name" json:"name"`
	Debt float64 `gorm:"column:credit" json:"credit"`
}

// FindUserUUIDsByName 根据客户名/手机号模糊查询 userUUID 列表
func (dao *customerDao) FindUserUUIDsByName(db *gorm.DB, ownerUser, nameSearch string) (uuids []string, err error) {
	like := "%" + nameSearch + "%"
	err = db.Model(&model.Customer{}).
		Where("owner_user = ? AND (name LIKE ? OR phone LIKE ?)", ownerUser, like, like).
		Pluck("uid", &uuids).Error
	return
}

// FindDebtCustomers 分页查询商户下有赊欠的客户（credit != 0，排除现金客户）
// search 非空时按 name / phone 模糊匹配
// 传 pageCount <= 0 表示不分页查全部
func (dao *customerDao) FindDebtCustomers(db *gorm.DB, ownerUser, search string, offset, limit int) (list []DebtCustomer, total int64, totalDebt float64, err error) {
	baseQuery := db.Model(&model.Customer{}).
		Where("owner_user = ? AND credit != 0 AND name <> ?", ownerUser, common.CashCustomerName)

	if search != "" {
		like := "%" + search + "%"
		baseQuery = baseQuery.Where("name LIKE ? OR phone LIKE ?", like, like)
	}

	// 总人数（用 Session 隔离，避免污染后续查询）
	if err = baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return
	}
	// 总欠款金额
	if err = baseQuery.Session(&gorm.Session{}).Select("COALESCE(SUM(credit), 0)").Scan(&totalDebt).Error; err != nil {
		return
	}

	// 分页列表
	query := baseQuery.Session(&gorm.Session{}).
		Select("uid, name, credit").
		Order("credit desc")
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}
	err = query.Scan(&list).Error
	return
}
