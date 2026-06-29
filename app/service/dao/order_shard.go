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

// ShardOrderDao 分表订单DAO
type ShardOrderDao struct{}

func newShardOrderDao() *ShardOrderDao {
	return &ShardOrderDao{}
}

// getShardTable 获取分表名
func (dao *ShardOrderDao) getShardTable(ownerUser string) string {
	return utils.GetShardTableName("batch_orders", ownerUser)
}

func (dao *ShardOrderDao) getShardGoodsTable(ownerUser string) string {
	return utils.GetShardTableName("batch_order_goods", ownerUser)
}

func (dao *ShardOrderDao) getShardPayTable(ownerUser string) string {
	return utils.GetShardTableName("batch_order_pay", ownerUser)
}

// Create 创建订单（自动路由到分表）
func (dao *ShardOrderDao) Create(db *gorm.DB, order *model.BatchOrder) error {
	tableName := dao.getShardTable(order.OwnerUser)
	return db.Table(tableName).Create(order).Error
}

// UpdateStatus 更新订单状态
func (dao *ShardOrderDao) UpdateStatus(db *gorm.DB, ownerUser, orderUUID string, status int) error {
	tableName := dao.getShardTable(ownerUser)
	return utils.WhereUIDCond(orderUUID).Cond(db).Table(tableName).Model(&model.BatchOrder{}).Update("status", status).Error
}

// Shared 分享订单
func (dao *ShardOrderDao) Shared(db *gorm.DB, ownerUser, orderUUID string) error {
	tableName := dao.getShardTable(ownerUser)
	return utils.WhereUIDCond(orderUUID).Cond(db).Table(tableName).Updates(&model.BatchOrder{
		Shared: common.BatchOrderShared,
	}).Error
}

// FindByUID 根据UID查询订单
func (dao *ShardOrderDao) FindByUID(db *gorm.DB, ownerUser, uid string, preload ...string) (*model.BatchOrder, error) {
	tableName := dao.getShardTable(ownerUser)
	query := db.Table(tableName).Where("uid = ?", uid)

	// 处理预加载
	for _, p := range preload {
		if p == "GoodsListRelated" {
			goodsTable := dao.getShardGoodsTable(ownerUser)
			query = query.Preload("GoodsListRelated", func(db *gorm.DB) *gorm.DB {
				return db.Table(goodsTable)
			})
		}
	}

	var order model.BatchOrder
	if err := query.First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// List 查询订单列表（分表）
func (dao *ShardOrderDao) List(db *gorm.DB, ownerUser string, userUUID string, startTime, endTime time.Time, status int32, offset, limit int, preloadGoods bool) ([]*model.BatchOrder, error) {
	tableName := dao.getShardTable(ownerUser)
	query := db.Table(tableName).Where("owner_user = ?", ownerUser)

	// 客户筛选
	if userUUID != "" {
		query = query.Where("user_uuid = ?", userUUID)
	}

	// 时间范围筛选
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	// 状态筛选
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	// 预加载货品
	if preloadGoods {
		goodsTable := dao.getShardGoodsTable(ownerUser)
		query = query.Preload("GoodsListRelated", func(db *gorm.DB) *gorm.DB {
			return db.Table(goodsTable)
		})
	}

	// 分页和排序
	var orderList []*model.BatchOrder
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orderList).Error
	return orderList, err
}

// CreateGoods 创建订单货品
func (dao *ShardOrderDao) CreateGoods(db *gorm.DB, goods *model.BatchOrderGoods) error {
	tableName := dao.getShardGoodsTable(goods.OwnerUser)
	return db.Table(tableName).Create(goods).Error
}

// CreateGoodsBatch 批量创建订单货品
func (dao *ShardOrderDao) CreateGoodsBatch(db *gorm.DB, ownerUser string, goodsList []*model.BatchOrderGoods) error {
	if len(goodsList) == 0 {
		return nil
	}
	tableName := dao.getShardGoodsTable(ownerUser)
	return db.Table(tableName).CreateInBatches(goodsList, 100).Error
}

// CreatePay 创建订单支付记录
func (dao *ShardOrderDao) CreatePay(db *gorm.DB, pay *model.BatchOrderPay) error {
	tableName := dao.getShardPayTable(pay.OwnerUser)
	return db.Table(tableName).Create(pay).Error
}

// DeleteGoodsByOrderUUID 根据订单UUID删除货品
func (dao *ShardOrderDao) DeleteGoodsByOrderUUID(db *gorm.DB, ownerUser, batchOrderUUID string) error {
	tableName := dao.getShardGoodsTable(ownerUser)
	return db.Table(tableName).Where("batch_order_uuid = ?", batchOrderUUID).Delete(&model.BatchOrderGoods{}).Error
}

// CustomerCreditAmount 实时查询单个客户的赊欠金额（分表）
func (dao *ShardOrderDao) CustomerCreditAmount(db *gorm.DB, ownerUser, userUUID string) (creditAmount float64, err error) {
	tableName := dao.getShardTable(ownerUser)
	err = db.Table(tableName).
		Where("owner_user = ? AND user_uuid = ?", ownerUser, userUUID).
		Where("status IN (?)", common.ValidOrder).
		Where("credit_amount > ?", 0).
		Select("COALESCE(SUM(credit_amount), 0)").
		Scan(&creditAmount).Error
	return
}

// CreditAmountTotal 统计赊欠总额（分表）
func (dao *ShardOrderDao) CreditAmountTotal(db *gorm.DB, ownerUser string) (creditAmount float64, err error) {
	tableName := dao.getShardTable(ownerUser)
	err = db.Table(tableName).
		Where("owner_user = ?", ownerUser).
		Where("status IN (?)", common.ValidOrder).
		Select("COALESCE(SUM(credit_amount), 0)").
		Scan(&creditAmount).Error

	if err != nil {
		return 0, fmt.Errorf("[CreditAmountTotal] 统计赊欠总额失败 %s: %w", ownerUser, err)
	}

	return creditAmount, nil
}

// MonthSales 月度销售额统计（分表）
func (dao *ShardOrderDao) MonthSales(db *gorm.DB, ownerUser string) (amount float64, err error) {
	tableName := dao.getShardTable(ownerUser)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	err = db.Table(tableName).
		Where("owner_user = ?", ownerUser).
		Where("created_at >= ?", monthStart).
		Where("status IN (?)", common.ValidOrder).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&amount).Error

	if err != nil {
		return 0, fmt.Errorf("[MonthSales] 统计月度金额失败 %s: %w", ownerUser, err)
	}

	return amount, nil
}

// LatestOrderByCustomers 查询客户最近订单时间（分表）
func (dao *ShardOrderDao) LatestOrderByCustomers(db *gorm.DB, ownerUser string, customers []model.Customer) (billingLatestDate map[string]int, err error) {
	billingLatestDate = make(map[string]int)
	customUUIDList := make([]string, 0)
	for _, customer := range customers {
		billingLatestDate[customer.UID] = utils.DurationDays(customer.CreatedAt)
		customUUIDList = append(customUUIDList, customer.UID)
	}

	tableName := dao.getShardGoodsTable(ownerUser)
	var results []struct {
		UserUUID      string
		LatestOrderAt time.Time
	}

	if err = db.Table(tableName).Where("user_uuid in (?)", customUUIDList).
		Select("user_uuid, MAX(created_at) as latest_order_at").
		Group("user_uuid").
		Scan(&results).Error; err != nil {
		return
	}

	for _, result := range results {
		billingLatestDate[result.UserUUID] = utils.DurationDays(result.LatestOrderAt)
	}

	return
}

// ListByBatchUUIDIn 根据批次UUID和货品UUID列表查询（分表）
func (dao *ShardOrderDao) ListByBatchUUIDIn(db *gorm.DB, ownerUser, batchUUID string, goodsUUIDList []string) (objects []model.BatchOrderGoods, err error) {
	tableName := dao.getShardGoodsTable(ownerUser)
	err = utils.Find(db.Table(tableName), &objects,
		utils.NewWhereCond("batch_uuid", batchUUID),
		utils.NewInCondFromString("goods_uuid", goodsUUIDList))
	return
}

// UpdateOrderPay 更新订单支付（分表）
func (dao *ShardOrderDao) UpdateOrderPay(db *gorm.DB, ownerUser string, object model.BatchOrderPay) error {
	tableName := dao.getShardPayTable(ownerUser)
	return utils.WhereUIDCond(object.UID).Cond(db).Table(tableName).Updates(&model.BatchOrderPay{
		PayType: object.PayType,
		Amount:  object.Amount,
	}).Error
}

// UpdateGoods 更新订单货品（分表）
func (dao *ShardOrderDao) UpdateGoods(db *gorm.DB, ownerUser string, object model.BatchOrderGoods) error {
	tableName := dao.getShardGoodsTable(ownerUser)
	return utils.WhereUIDCond(object.UID).Cond(db).Table(tableName).Updates(&model.BatchOrderGoods{
		Price:  object.Price,
		Mount:  object.Mount,
		Weight: object.Weight,
	}).Error
}

// GetTodayOrdersWithGoods 获取指定订单或当日的订单及货品（分表）
// 如果 orderUUID 不为空，则查询该订单所在日期的数据
// 如果 orderUUID 为空，则查询当日的数据
func (dao *ShardOrderDao) GetTodayOrdersWithGoods(db *gorm.DB, ownerUser, customerUUID, orderUUID string) (orders []model.BatchOrder, previousCredit float64, err error) {
	orderTable := dao.getShardTable(ownerUser)
	goodsTable := dao.getShardGoodsTable(ownerUser)

	var targetDate time.Time
	var dayStart, dayEnd time.Time

	if orderUUID != "" {
		// 1. 先查询订单的创建时间
		var order model.BatchOrder
		err = db.Table(orderTable).Where("uid = ?", orderUUID).First(&order).Error
		if err != nil {
			return nil, 0, fmt.Errorf("订单不存在: %w", err)
		}
		targetDate = order.CreatedAt
	} else {
		// 使用当前日期
		targetDate = time.Now()
	}

	// 计算目标日期的开始和结束时间
	dayStart = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	dayEnd = dayStart.Add(24 * time.Hour)

	// 1. Sum of credit amounts for valid orders BEFORE the target day
	err = db.Table(orderTable).
		Where("owner_user = ? AND user_uuid = ? AND status IN (?) AND credit_amount > 0 AND created_at < ?",
			ownerUser, customerUUID, common.ValidOrder, dayStart).
		Select("COALESCE(SUM(credit_amount), 0)").
		Scan(&previousCredit).Error
	if err != nil {
		return
	}

	// 2. Fetch the target day's orders with preloaded goods
	err = db.Table(orderTable).
		Where("owner_user = ? AND user_uuid = ? AND status IN (?) AND created_at >= ? AND created_at < ?",
			ownerUser, customerUUID, common.ValidOrder, dayStart, dayEnd).
		Preload("GoodsListRelated", func(db *gorm.DB) *gorm.DB {
			return db.Table(goodsTable)
		}).
		Find(&orders).Error

	return
}

// GetCreditList 赊欠列表查询（分表）
func (dao *ShardOrderDao) GetCreditList(db *gorm.DB, ownerUser string, conditions request.CreditListReq) (*response.CreditListResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tableName := dao.getShardTable(ownerUser)

	// 基础查询条件复用
	baseQuery := db.Table(tableName).
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
