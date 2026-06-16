package dao

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type orderDao struct{}

func newOrderDao() *orderDao {
	return &orderDao{}
}

func (dao *orderDao) UpdateStatus(db *gorm.DB, orderUUID string, status int) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Model(&model.BatchOrder{}).Update("status", status).Error
}

func (dao *orderDao) Shared(db *gorm.DB, orderUUID string) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Updates(&model.BatchOrder{
		Shared: common.BatchOrderShared,
	}).Error
}

func (dao *orderDao) UpdateGoods(db *gorm.DB, object model.BatchOrderGoods) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchOrderGoods{
		Price:  object.Price,
		Mount:  object.Mount,
		Weight: object.Weight,
	}).Error
}

func (dao *orderDao) UpdateOrderPay(db *gorm.DB, object model.BatchOrderPay) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchOrderPay{
		PayType: object.PayType,
		Amount:  object.Amount,
	}).Error
}

func (dao *orderDao) LatestOrderByCustomers(db *gorm.DB, owneruser string, customers []model.Customer) (billingLatestDate map[string]int, err error) {
	billingLatestDate = make(map[string]int)
	customUUIDList := make([]string, 0)
	for _, customer := range customers {
		billingLatestDate[customer.UID] = utils.DurationDays(customer.CreatedAt)
		customUUIDList = append(customUUIDList, customer.UID)
	}
	var results []struct {
		UserUUID      string
		LatestOrderAt time.Time
	}

	if err = db.Model(&model.BatchOrderGoods{}).Where("user_uuid in (?)", customUUIDList).
		Select("user_uuid, MAX(created_at) as latest_order_at").
		Group("user_uuid").
		Scan(&results).Error; err != nil {
		global.Global.Logger.Error(fmt.Sprintf("[BillingCondByOwnerUser] %s %s", owneruser, err))
		return
	}

	for _, result := range results {
		billingLatestDate[result.UserUUID] = utils.DurationDays(result.LatestOrderAt)
	}

	return
}

func (dao *orderDao) MonthSales(db *gorm.DB, ownerUser string) (amount float64, err error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 使用数据库聚合查询直接求和
	err = db.Model(&model.BatchOrder{}).
		Where("owner_user = ?", ownerUser).
		Where("created_at >= ?", monthStart).
		Where("status IN (?)", common.ValidOrder).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&amount).Error

	if err != nil {
		global.Global.Logger.Error("[MonthFinance] 统计月度金额失败",
			zap.String("ownerUser", ownerUser),
			zap.Error(err))
		return 0, err
	}

	global.Global.Logger.Debug("[MonthFinance] 月度统计",
		zap.String("ownerUser", ownerUser),
		zap.Float64("totalAmount", amount))

	return amount, nil
}

func (dao *orderDao) ListByBatchUUIDIn(db *gorm.DB, batchUUID string, goodsUUIDList []string) (objects []model.BatchOrderGoods, err error) {
	err = utils.Find(db, objects, utils.NewWhereCond("batch_uuid", batchUUID), utils.NewInCondFromString("goods_uuid", goodsUUIDList))
	return
}

// 仅统计赊欠总额度（使用聚合查询，效率更高）
func (dao *orderDao) CreditAmountTotal(db *gorm.DB, ownerUser string) (creditAmount float64, err error) {
	err = db.Model(&model.BatchOrder{}).
		Where("owner_user = ?", ownerUser).
		Where("status IN (?)", common.ValidOrder).
		Select("COALESCE(SUM(credit_amount), 0)").
		Scan(&creditAmount).Error

	if err != nil {
		global.Global.Logger.Error(fmt.Sprintf("[CreditAmountTotal] 统计赊欠总额失败 %s: %s", ownerUser, err))
		return 0, err
	}

	return creditAmount, nil
}

// CustomerCreditAmount 实时查询单个客户的赊欠金额
func (dao *orderDao) CustomerCreditAmount(db *gorm.DB, ownerUser, userUUID string) (creditAmount float64, err error) {
	err = db.Model(&model.BatchOrder{}).
		Where("owner_user = ? AND user_uuid = ?", ownerUser, userUUID).
		Where("status IN (?)", common.ValidOrder).
		Where("credit_amount > ?", 0).
		Select("COALESCE(SUM(credit_amount), 0)").
		Scan(&creditAmount).Error
	return
}

// GetCreditList 优化后的赊欠列表查询
func (dao *orderDao) GetCreditList(db *gorm.DB, ownerUser string, conditions request.CreditListReq) (*response.CreditListResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 基础查询条件复用
	baseQuery := db.Table("batch_orders").
		Where("owner_user = ?", ownerUser).
		Where("status = ?", common.BatchOrderTemp).
		Where("credit_amount > ?", 0)

	// 1. 先查询真实的汇总数据（全量统计）
	var summary response.CreditSummary
	err := baseQuery.
		Select(`COALESCE(SUM(credit_amount), 0) as total_credit_amount,
			COUNT(DISTINCT user_uuid) as total_credit_users,
			COUNT(*) as total_credit_orders,
			COALESCE(SUM(CASE WHEN created_at >= ? THEN credit_amount ELSE 0 END), 0) as today_credit_amount
		`, todayStart).
		Scan(&summary).Error

	if err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询赊欠汇总失败", zap.Error(err))
		return nil, err
	}

	// 2. 如果有赊欠数据才查询列表
	if summary.TotalCreditUsers == 0 {
		return &response.CreditListResponse{
			List:    []*response.CreditUserStat{},
			Summary: summary,
		}, nil
	}

	// 3. 查询分页列表（每个user_uuid最早的赊欠时间）
	type earliestCredit struct {
		UserUUID     string    `gorm:"column:user_uuid"`
		FirstCreated time.Time `gorm:"column:first_created"`
	}

	// 子查询：获取每个用户最早的赊欠时间
	earliestSub := baseQuery.
		Select("user_uuid, MIN(created_at) as first_created").
		Group("user_uuid")

	var earliestList []earliestCredit
	err = db.Table("(?) as ec", earliestSub).Scan(&earliestList).Error
	if err != nil {
		return nil, err
	}

	global.Global.Logger.Debug("earliestList", zap.Any("earliestList", earliestList))

	// 构建最早时间映射（减少JOIN复杂度）
	earliestMap := make(map[string]time.Time, len(earliestList))
	for _, ec := range earliestList {
		earliestMap[ec.UserUUID] = ec.FirstCreated
	}

	// 4. 查询分组统计（支持分页）
	type creditStat struct {
		UserUUID    string  `gorm:"column:user_uuid"`
		TotalCredit float64 `gorm:"column:total_credit"`
		TodayCredit float64 `gorm:"column:today_credit"`
		OrderCount  int64   `gorm:"column:order_count"`
	}

	var stats []creditStat
	statQuery := baseQuery.
		Select(`
			user_uuid,
			SUM(credit_amount) as total_credit,
			SUM(CASE WHEN created_at >= ? THEN credit_amount ELSE 0 END) as today_credit,
			COUNT(*) as order_count
		`, todayStart).
		Group("user_uuid").
		Order("total_credit DESC")

	// 分页处理
	if !conditions.LoadAll && conditions.PageCount > 0 {
		offset := (conditions.Page - 1) * conditions.PageCount
		statQuery = statQuery.Offset(offset).Limit(conditions.PageCount)
	}

	err = statQuery.Scan(&stats).Error
	if err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询赊欠列表失败", zap.Error(err))
		return nil, err
	}

	// 5. 构建列表（结合最早时间计算赊欠天数）
	list := make([]*response.CreditUserStat, len(stats))
	for i, s := range stats {
		longestDays := 0
		if firstTime, ok := earliestMap[s.UserUUID]; ok {
			longestDays = int(now.Sub(firstTime).Hours() / 24)
		}

		list[i] = &response.CreditUserStat{
			UserUUID:    s.UserUUID,
			TotalCredit: s.TotalCredit,
			TodayCredit: s.TodayCredit,
			LongestDays: longestDays,
			OrderCount:  s.OrderCount,
		}
	}

	global.Global.Logger.Info("GetCreditList", zap.Any("list", list))
	return &response.CreditListResponse{
		List:    list,
		Summary: summary,
	}, nil
}

// GetTodayOrdersWithGoods retrieves today's BatchOrders (status=1) for a given owner and customer,
// along with their BatchOrderGoods. Returns the order list and the sum of CreditAmount before today.
func (dao *orderDao) GetTodayOrdersWithGoods(db *gorm.DB, ownerUser, customerUUID string) (*response.ShareDailyOrderRsp, error) {
	// now := time.Now()
	// todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// todayEnd := todayStart.Add(24 * time.Hour)

	var resp response.ShareDailyOrderRsp

	// // 1. Calculate previous credit sum (all status=1 orders for this customer before today)
	// err := db.Model(&model.BatchOrder{}).
	// 	Where("owner_user = ? AND user_uuid = ? AND status = ? AND credit_amount > 0 AND created_at < ?",
	// 		ownerUser, customerUUID, 1, todayStart).
	// 	Select("COALESCE(SUM(credit_amount), 0)").
	// 	Scan(&resp.TotalPreviousCredit).Error
	// if err != nil {
	// 	global.Global.Logger.Error("failed to get previous credit", zap.Error(err))
	// 	return nil, err
	// }

	// // 2. Fetch today's orders with status=1
	// var orders []model.BatchOrder
	// err = db.Where("owner_user = ? AND user_uuid = ? AND status = ? AND created_at >= ? AND created_at < ?",
	// 	ownerUser, customerUUID, 1, todayStart, todayEnd).
	// 	Preload("GoodsListRelated"). // Preload associated BatchOrderGoods
	// 	Find(&orders).Error
	// if err != nil {
	// 	global.Global.Logger.Error("failed to get today orders", zap.Error(err))
	// 	return nil, err
	// }

	// shareOrder := &response.ShareDailyOrderRsp{
	// 	TotalAmount:  order.TotalAmount,
	// 	CreditAmount: order.CreditAmount,
	// 	GoodsList:    make([]*response.ShareDailyOrderItem, 0, len(order.GoodsListRelated)),
	// }
	// // 3. Build response
	// for _, order := range orders {

	// 	for _, good := range order.GoodsListRelated {
	// 		item := &response.ShareDailyOrderItem{
	// 			GoodsUUID: good.GoodsUUID,
	// 			Price:     good.Price,
	// 			Weight:    good.Weight,
	// 			Mount:     good.Mount,
	// 			Total:     good.Total,
	// 		}
	// 		shareOrder.GoodsList = append(shareOrder.GoodsList, item)
	// 	}
	// }

	return &resp, nil
}
