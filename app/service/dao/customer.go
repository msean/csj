package dao

import (
	"app/pkg/utils"
	"app/service/model"
	"app/service/model/request"

	"gorm.io/gorm"
)

type customerDao struct{}

func newCustomerDao() *customerDao {
	return &customerDao{}
}

func (dao *customerDao) Update(db *gorm.DB, object model.Customer) (err error) {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.Customer{
		Name:   object.Name,
		Phone:  object.Phone,
		Remark: object.Remark,
		Debt:   object.Debt,
		Addr:   object.Addr,
	}).Error
}

func (dao *customerDao) NewTempCustomer(db *gorm.DB, ownerUser string) (err error) {
	c := model.Customer{
		Name:      "现金客户",
		Remark:    "现金客户",
		OwnerUser: ownerUser,
	}
	return utils.CreateObj(db, &c)
}

func (dao *customerDao) List(db *gorm.DB, ownerUser string, condtions request.CustomerListReq) (customers []model.Customer, err error) {
	conds := []utils.Cond{}
	conds = append(conds, []utils.Cond{
		utils.NewWhereCond("owner_user", ownerUser),
		utils.NewOrderCond("Convert(name USING gbk)"),
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
func (dao *customerDao) IncrOrderStats(db *gorm.DB, userUUID string, totalAmount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Updates(map[string]interface{}{
			"order_count":  gorm.Expr("order_count + 1"),
			"total_amount": gorm.Expr("total_amount + ?", totalAmount),
		}).Error
	return
}

// DecrOrderStats 作废/退款/退货时减少客户统计
func (dao *customerDao) DecrOrderStats(db *gorm.DB, userUUID string, totalAmount float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Updates(map[string]interface{}{
			"order_count":  gorm.Expr("GREATEST(order_count - 1, 0)"),
			"total_amount": gorm.Expr("total_amount - ?", totalAmount),
		}).Error
	return
}

// UpdateOrderStatsDiff 订单更新时调整客户统计差额
func (dao *customerDao) UpdateOrderStatsDiff(db *gorm.DB, userUUID string, deltaTotal float64) (err error) {
	err = db.Model(&model.Customer{}).
		Where("uid = ?", userUUID).
		Update("total_amount", gorm.Expr("total_amount + ?", deltaTotal)).Error
	return
}
