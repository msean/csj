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

// orderDao 分表订单DAO
type orderDao struct{}

func newOrderDao() *orderDao {
	return &orderDao{}
}

// Create 创建订单（自动路由到分表）
func (dao *orderDao) Create(db *gorm.DB, order *model.BatchOrder) error {
	return db.Table(ShardOrderTable(order.OwnerUser)).Create(order).Error
}

// UpdateStatus 更新订单状态
func (dao *orderDao) UpdateStatus(db *gorm.DB, ownerUser, orderUUID string, status int) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Table(ShardOrderTable(ownerUser)).Model(&model.BatchOrder{}).Update("status", status).Error
}

// FindByUID 根据UID查询订单
func (dao *orderDao) FindByUID(db *gorm.DB, ownerUser, uid string, preload ...string) (order model.BatchOrder, err error) {
	// 处理预加载
	for _, p := range preload {
		if p == "GoodsListRelated" {
			db = db.Preload("GoodsListRelated", func(_db *gorm.DB) *gorm.DB {
				return _db.Table(ShardOrderGoodsTable(ownerUser))
			})
		}
	}

	err = utils.First(db.Table(ShardOrderTable(ownerUser)), &order, utils.WhereUIDCond(uid))
	return order, nil
}

// List 查询订单列表（分表）
func (dao *orderDao) List(db *gorm.DB, ownerUser string, userUUIDs []string, startTime, endTime time.Time, status int32, offset, limit int, preloadGoods bool) ([]*model.BatchOrder, error) {
	query := db.Table(ShardOrderTable(ownerUser)).Where("owner_user = ?", ownerUser)

	// 客户筛选
	if len(userUUIDs) > 0 {
		query = query.Where("user_uuid IN (?)", userUUIDs)
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
		query = query.Preload("GoodsListRelated", func(db *gorm.DB) *gorm.DB {
			return db.Table(ShardOrderGoodsTable(ownerUser))
		})
	}

	// 分页和排序
	var orderList []*model.BatchOrder
	err := query.Order("order_time DESC").Offset(offset).Limit(limit).Find(&orderList).Error
	return orderList, err
}

// CreateGoods 创建订单货品
func (dao *orderDao) CreateGoods(db *gorm.DB, goods *model.BatchOrderGoods) error {
	return db.Table(ShardOrderGoodsTable(goods.OwnerUser)).Create(goods).Error
}

// CreateGoodsBatch 批量创建订单货品
func (dao *orderDao) CreateGoodsBatch(db *gorm.DB, ownerUser string, goodsList []*model.BatchOrderGoods) error {
	if len(goodsList) == 0 {
		return nil
	}
	return db.Table(ShardOrderGoodsTable(ownerUser)).CreateInBatches(goodsList, 100).Error
}

// CreatePay 创建订单支付记录
func (dao *orderDao) CreatePay(db *gorm.DB, pay *model.BatchOrderPay) error {
	return db.Table(ShardPayTable(pay.OwnerUser)).Create(pay).Error
}

// DeleteGoodsByOrderUUID 根据订单UUID删除货品
func (dao *orderDao) DeleteGoodsByOrderUUID(db *gorm.DB, ownerUser, batchOrderUUID string) error {
	return db.Table(ShardOrderGoodsTable(ownerUser)).Where("batch_order_uuid = ?", batchOrderUUID).Delete(&model.BatchOrderGoods{}).Error
}

// DeletePayByOrderUUID 根据订单UUID删除还款记录（分表）
func (dao *orderDao) DeletePayByOrderUUID(db *gorm.DB, ownerUser, batchOrderUUID string) error {
	return db.Table(ShardPayTable(ownerUser)).Where("batch_order_uuid = ?", batchOrderUUID).Delete(&model.BatchOrderPay{}).Error
}

// UpdateOrderFields 更新订单指定字段（分表）
// 仅更新传入的字段，避免 GORM 处理 GoodsListRelated 关联触犯外键约束。
func (dao *orderDao) UpdateOrderFields(db *gorm.DB, ownerUser, uid string, fields map[string]interface{}) error {
	return db.Table(ShardOrderTable(ownerUser)).Where("uid = ?", uid).Updates(fields).Error
}

// UpdateOrderCreditStatus 更新订单的赊欠金额和状态（分表）
// 用于还款后/撤销还款后调整订单状态
func (dao *orderDao) UpdateOrderCreditStatus(db *gorm.DB, ownerUser, uid string, creditAmount float64, status int) error {
	return db.Table(ShardOrderTable(ownerUser)).Where("uid = ?", uid).
		Updates(map[string]interface{}{
			"credit_amount": creditAmount,
			"status":        status,
		}).Error
}

// CustomerCreditAmount 实时查询单个客户的赊欠金额（分表）
func (dao *orderDao) CustomerCreditAmount(db *gorm.DB, ownerUser, userUUID string) (creditAmount float64, err error) {
	err = db.Table(ShardOrderTable(ownerUser)).
		Where("owner_user = ? AND user_uuid = ?", ownerUser, userUUID).
		Where("status IN (?)", common.UnFinalStatusOrder).
		Where("credit_amount > ?", 0).
		Select("COALESCE(SUM(credit_amount), 0)").
		Scan(&creditAmount).Error
	return
}

func (dao *orderDao) FindOrdersWithCredit(db *gorm.DB, ownerUser, customerUUID string, validStatuses []int) ([]model.BatchOrder, error) {
	var orders []model.BatchOrder
	err := db.Table(ShardOrderTable(ownerUser)).Where("owner_user = ? AND user_uuid = ? AND status IN (?) AND credit_amount != 0",
		ownerUser, customerUUID, validStatuses).
		Order("created_at ASC"). // 最早的订单优先
		Find(&orders).Error
	return orders, err
}

// CreditAmountTotal 统计赊欠总额（从 Customer 表汇总所有客户的 Credit）
func (dao *orderDao) CreditAmountTotal(db *gorm.DB, ownerUser string) (creditAmount float64, err error) {
	err = db.Model(&model.Customer{}).
		Where("owner_user = ?", ownerUser).
		Select("COALESCE(SUM(credit), 0)").
		Scan(&creditAmount).Error

	if err != nil {
		return 0, fmt.Errorf("[CreditAmountTotal] 统计赊欠总额失败 %s: %w", ownerUser, err)
	}

	return creditAmount, nil
}

// FindCreditSummariesWithPage 分页查询有赊欠的客户汇总
func (dao *orderDao) FindCreditSummariesWithPage(db *gorm.DB, ownerUser string, validStatuses []int, offset, limit int) ([]CreditSummary, int64, error) {
	var summaries []CreditSummary
	var total int64

	query := db.Table(ShardOrderTable(ownerUser)).
		Select("user_uuid, SUM(credit_amount) as total_credit").
		Where("owner_user = ? AND status IN (?) AND credit_amount != 0", ownerUser, validStatuses).
		Group("user_uuid")

	// 获取总数
	query.Count(&total)

	// 分页
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Scan(&summaries).Error
	return summaries, total, err
}

// FindCreditSummaries 查询所有有赊欠的客户汇总
func (dao *orderDao) FindCreditSummaries(db *gorm.DB, ownerUser string, validStatuses []int) ([]CreditSummary, error) {
	var summaries []CreditSummary
	err := db.Table(ShardOrderTable(ownerUser)).
		Select("user_uuid, SUM(credit_amount) as total_credit").
		Where("owner_user = ? AND status IN (?) AND credit_amount != 0", ownerUser, validStatuses).
		Group("user_uuid").
		Scan(&summaries).Error
	return summaries, err
}

// MonthSales 月度销售额统计（分表）
func (dao *orderDao) MonthSales(db *gorm.DB, ownerUser string) (amount float64, err error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err = db.Table(ShardOrderTable(ownerUser)).
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
func (dao *orderDao) LatestOrderByCustomers(db *gorm.DB, ownerUser string, customers []model.Customer) (billingLatestDate map[string]int, err error) {
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

	if err = db.Table(ShardOrderGoodsTable(ownerUser)).Where("user_uuid in (?)", customUUIDList).
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
func (dao *orderDao) ListByBatchUUIDIn(db *gorm.DB, ownerUser, batchUUID string, goodsUUIDList []string) (objects []model.BatchOrderGoods, err error) {
	err = utils.Find(db.Table(ShardOrderGoodsTable(ownerUser)), &objects,
		utils.NewWhereCond("batch_uuid", batchUUID),
		utils.NewInCondFromString("goods_uuid", goodsUUIDList))
	return
}

// UpdateGoods 更新订单货品（分表）
func (dao *orderDao) UpdateGoods(db *gorm.DB, ownerUser string, object model.BatchOrderGoods) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Table(ShardOrderGoodsTable(ownerUser)).Updates(&model.BatchOrderGoods{
		Price:  object.Price,
		Mount:  object.Mount,
		Weight: object.Weight,
	}).Error
}

// FindGoodsByBatchAndGoodsUUID 按批次UUID和货品UUID查询订单货品列表（分表）
// statuses 非空时，JOIN 订单表按订单状态筛选；按创建时间倒序排列，用于商品销售明细
func (dao *orderDao) FindGoodsByBatchAndGoodsUUID(db *gorm.DB, ownerUser, batchUUID, goodsUUID string, statuses []int) (list []model.BatchOrderGoods, err error) {
	goodsTable := ShardOrderGoodsTable(ownerUser)
	query := db.Table(goodsTable).
		Where(goodsTable+".batch_uuid = ? AND "+goodsTable+".goods_uuid = ?", batchUUID, goodsUUID)

	if len(statuses) > 0 {
		orderTable := ShardOrderTable(ownerUser)
		query = query.Joins("JOIN "+orderTable+" ON "+orderTable+".uid = "+goodsTable+".batch_order_uuid").
			Where(orderTable+".status IN (?)", statuses)
	}

	err = query.Order(goodsTable + ".created_at DESC").Find(&list).Error
	return
}

// GetShareOrdersWithGoods 获取分享所需的订单及货品（分表）
// 如果 orderUUID 不为空，返回该指定订单；否则返回当天或昨天最近的一单
func (dao *orderDao) GetShareOrdersWithGoods(db *gorm.DB, ownerUser, customerUUID, orderUUID string) (orders []model.BatchOrder, previousCredit float64, err error) {
	orderTable := ShardOrderTable(ownerUser)
	goodsTable := ShardOrderGoodsTable(ownerUser)

	var targetOrder model.BatchOrder
	if orderUUID != "" {
		err = db.Table(orderTable).
			Where("uid = ? AND owner_user = ? AND user_uuid = ? AND status IN (?)",
				orderUUID, ownerUser, customerUUID, common.UnFinalStatusOrder).
			First(&targetOrder).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = fmt.Errorf("订单不存在")
			}
			return
		}
	} else {
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		yesterdayStart := todayStart.Add(-24 * time.Hour)
		todayEnd := todayStart.Add(24 * time.Hour)

		err = db.Table(orderTable).
			Where("owner_user = ? AND user_uuid = ? AND status IN (?) AND order_time >= ? AND order_time < ?",
				ownerUser, customerUUID, common.UnFinalStatusOrder, yesterdayStart.Unix(), todayEnd.Unix()).
			Order("order_time DESC").
			First(&targetOrder).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
			}
			return
		}
	}

	err = db.Table(orderTable).
		Where("uid = ?", targetOrder.UID).
		Preload("GoodsListRelated", func(db *gorm.DB) *gorm.DB {
			return db.Table(goodsTable)
		}).
		Find(&orders).Error
	return
}

// GetCreditList 赊欠列表查询
// 总赊欠额从 Customer 表查询（包含客户初始设置的赊欠余额），
// 今日赊欠、订单数、最长赊欠天数仍从订单分表查询
func (dao *orderDao) GetCreditList(db *gorm.DB, ownerUser string, conditions request.CreditListReq) (*response.CreditListResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 1. 从 Customer 表查询汇总（总赊欠额、赊欠人数）
	var summary response.CreditSummary
	err := db.Model(&model.Customer{}).
		Where("owner_user = ? AND credit != 0 AND name <> ?", ownerUser, common.CashCustomerName).
		Select("COALESCE(SUM(credit), 0) as total_credit_amount, COUNT(*) as total_credit_users").
		Scan(&summary).Error
	if err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询客户赊欠汇总失败", zap.Error(err))
		return nil, err
	}

	// 2. 从订单分表查询订单维度的汇总（赊欠总单数、今日赊欠总额）
	orderTable := ShardOrderTable(ownerUser)
	var orderSummary struct {
		TotalCreditOrders int64   `gorm:"column:total_credit_orders"`
		TodayCreditAmount float64 `gorm:"column:today_credit_amount"`
	}
	err = db.Table(orderTable+" AS o").
		Where("o.owner_user = ? AND o.status = ? AND o.credit_amount > 0", ownerUser, common.BatchOrderTemp).
		Joins("JOIN customers c ON c.uid = o.user_uuid").
		Where("c.name <> ?", common.CashCustomerName).
		Select(`COUNT(*) as total_credit_orders,
			COALESCE(SUM(CASE WHEN o.order_time >= ? THEN o.credit_amount ELSE 0 END), 0) as today_credit_amount`, todayStart.Unix()).
		Scan(&orderSummary).Error
	if err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询订单赊欠汇总失败", zap.Error(err))
		return nil, err
	}
	summary.TotalCreditOrders = orderSummary.TotalCreditOrders
	summary.TodayCreditAmount = orderSummary.TodayCreditAmount

	// 3. 无赊欠客户直接返回
	if summary.TotalCreditUsers == 0 {
		return &response.CreditListResponse{
			List:    []*response.CreditUserStat{},
			Summary: summary,
		}, nil
	}

	// 4. 从 Customer 表分页查询有赊欠的客户列表
	type customerCredit struct {
		UID    string  `gorm:"column:uid"`
		Name   string  `gorm:"column:name"`
		Credit float64 `gorm:"column:credit"`
	}

	customerQuery := db.Model(&model.Customer{}).
		Select("uid, name, credit").
		Where("owner_user = ? AND credit != 0 AND name <> ?", ownerUser, common.CashCustomerName)

	if conditions.Search != "" {
		like := "%" + conditions.Search + "%"
		customerQuery = customerQuery.Where("name LIKE ? OR phone LIKE ?", like, like)
	}

	customerQuery = customerQuery.Order("credit DESC")

	if !conditions.LoadAll && conditions.PageCount > 0 {
		offset := (conditions.Page - 1) * conditions.PageCount
		customerQuery = customerQuery.Offset(offset).Limit(conditions.PageCount)
	}

	var customers []customerCredit
	if err = customerQuery.Scan(&customers).Error; err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询赊欠客户列表失败", zap.Error(err))
		return nil, err
	}

	if len(customers) == 0 {
		return &response.CreditListResponse{
			List:    []*response.CreditUserStat{},
			Summary: summary,
		}, nil
	}

	// 5. 批量查询这些客户的订单维度数据（今日赊欠、订单数、最早赊欠时间）
	userUUIDs := make([]string, len(customers))
	for i, c := range customers {
		userUUIDs[i] = c.UID
	}

	type orderStat struct {
		UserUUID     string    `gorm:"column:user_uuid"`
		TodayCredit  float64   `gorm:"column:today_credit"`
		OrderCount   int64     `gorm:"column:order_count"`
		FirstCreated time.Time `gorm:"column:first_created"`
	}

	var orderStats []orderStat
	err = db.Table(orderTable+" AS o").
		Where("o.owner_user = ? AND o.user_uuid IN ? AND o.status = ? AND o.credit_amount > 0",
			ownerUser, userUUIDs, common.BatchOrderTemp).
		Select(`o.user_uuid,
			SUM(CASE WHEN o.order_time >= ? THEN o.credit_amount ELSE 0 END) as today_credit,
			COUNT(*) as order_count,
			MIN(o.created_at) as first_created`, todayStart.Unix()).
		Group("o.user_uuid").
		Scan(&orderStats).Error
	if err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询客户订单统计失败", zap.Error(err))
		return nil, err
	}

	orderStatMap := make(map[string]orderStat, len(orderStats))
	for _, s := range orderStats {
		orderStatMap[s.UserUUID] = s
	}

	// 6. 查询当天已分享的客户集合
	sharedSet, err := ShareDao.SharedTodayCustomers(db, ownerUser)
	if err != nil {
		global.Global.Logger.Error("[GetCreditList] 查询当日分享记录失败", zap.Error(err))
		return nil, err
	}

	// 7. 组装结果
	list := make([]*response.CreditUserStat, len(customers))
	for i, c := range customers {
		longestDays := 0
		todayCredit := 0.0
		var orderCount int64
		if s, ok := orderStatMap[c.UID]; ok {
			todayCredit = s.TodayCredit
			orderCount = s.OrderCount
			if !s.FirstCreated.IsZero() {
				longestDays = int(now.Sub(s.FirstCreated).Hours() / 24)
			}
		}

		list[i] = &response.CreditUserStat{
			UserUUID:    c.UID,
			UserName:    c.Name,
			TotalCredit: c.Credit,
			TodayCredit: todayCredit,
			LongestDays: longestDays,
			OrderCount:  orderCount,
			TodayShared: sharedSet[c.UID],
		}
	}

	return &response.CreditListResponse{
		List:    list,
		Summary: summary,
	}, nil
}
