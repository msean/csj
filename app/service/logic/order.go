package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/cache"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	OrderLogic struct {
		context *gin.Context
		runtime *global.RunTime
		model.BatchOrder
		History model.BatchOrderHistory `json:"history"`
	}
	OrderGoodsLogic struct {
		context *gin.Context
		runtime *global.RunTime
		*model.BatchOrderGoods
	}
)

func NewOrderLogic(context *gin.Context) *OrderLogic {
	logic := &OrderLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func NewOrderGoodsLogic(context *gin.Context) *OrderGoodsLogic {
	logic := &OrderGoodsLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

// isCashCustomer 判断当前订单的客户是否为现金客户
func (logic *OrderLogic) isCashCustomer() (bool, error) {
	if logic.UserUUID == "" {
		return false, nil
	}
	customer, err := dao.CustomerDao.FromUUID(logic.runtime.DB, logic.UserUUID)
	if err != nil {
		return false, err
	}
	return customer.Name == common.CashCustomerName, nil
}

// orderDedupTTL 同一客户在该时间段内提交相同货品时只允许创建一次订单
const orderDedupTTL = 3 * time.Second

// goodsFingerprint 生成货品明细指纹，用于识别两次提交是否为同一笔订单。
func goodsFingerprint(goodsList []*model.BatchOrderGoods) string {
	parts := make([]string, 0, len(goodsList))
	for _, goods := range goodsList {
		parts = append(parts, fmt.Sprintf("%s|%d|%.2f|%d|%.2f",
			goods.GoodsUUID, goods.GoodType, goods.Price, goods.Mount, goods.Weight))
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

// orderDedupKey 生成 Redis 幂等键：同一用户、同一客户、相同货品时键相同。
func (logic *OrderLogic) orderDedupKey() string {
	sum := sha256.Sum256([]byte(logic.UserUUID + "|" + goodsFingerprint(logic.GoodsListRelated)))
	return fmt.Sprintf("order:dedup:%s:%x", logic.OwnerUser, sum)
}

// acquireOrderDedup 用 Redis SETNX 抢占本次下单的幂等键。
// 键已存在说明 2 秒内已提交过相同订单，直接返回 false。
// Redis 异常时不阻断正常下单，退化为不做防重。
func (logic *OrderLogic) acquireOrderDedup() bool {
	if global.Global.Redis == nil {
		return true
	}
	ok, err := global.Global.Redis.SetNX(context.Background(), logic.orderDedupKey(), 1, orderDedupTTL).Result()
	if err != nil {
		global.Global.Logger.Warn("[OrderCreate] Redis 防重失败，跳过防重", zap.Error(err))
		return true
	}
	return ok
}

// 码单
func (logic *OrderLogic) TempCreate() (err error) {
	logic.BatchOrder.DefaultSet()

	// 前端未传下单时间时用当前时间
	if logic.OrderTime == 0 {
		logic.OrderTime = int(time.Now().Unix())
	}

	if logic.BatchUUID == "" {
		return common.BatchUUIDRequireErr
	}
	if len(logic.GoodsListRelated) == 0 {
		return common.BatchOrderGoodsRequireErr
	}

	// 在创建前加这一行
	if err = logic.FillGoodsWeightWithTotal(logic.runtime.DB); err != nil {
		return
	}

	// 现金客户不能记账单（付款金额必须 >= 订单金额，暂存单不接受付款故直接拦截）
	isCash, err := logic.isCashCustomer()
	if err != nil {
		return
	}
	if isCash {
		return common.CashCustomerPayLessErr
	}

	// 同一客户 2 秒内重复提交相同货品时直接拦截
	if !logic.acquireOrderDedup() {
		return common.OrderDuplicateSubmitErr
	}

	// 创建订单 + 更新客户统计（使用事务保证一致性）
	tx := logic.runtime.DB.Begin()

	// 1. 先创建订单（不包含货品）
	goodsList := logic.GoodsListRelated
	logic.GoodsListRelated = nil // 暂时清空，避免GORM通过外键创建

	if err = dao.OrderDao.Create(tx, &logic.BatchOrder); err != nil {
		logic.GoodsListRelated = goodsList // 恢复
		tx.Rollback()
		global.Global.Logger.Error("OrderLogic TempCreate CreateObj", zap.Any("logic.BatchOrder", logic.BatchOrder), zap.Error(err))
		return
	}

	// 2. 再创建订单货品
	logic.GoodsListRelated = goodsList // 恢复
	for _, goods := range logic.GoodsListRelated {
		goods.BatchOrderUID = logic.UID
		goods.OwnerUser = logic.OwnerUser
		goods.UserUUID = logic.UserUUID
		goods.BatchUUID = logic.BatchUUID
	}

	if err = dao.OrderDao.CreateGoodsBatch(tx, logic.OwnerUser, logic.GoodsListRelated); err != nil {
		tx.Rollback()
		global.Global.Logger.Error("OrderLogic TempCreate CreateGoodsBatch", zap.Error(err))
		return
	}

	if err = dao.CustomerDao.IncrOrderStats(tx, logic.UserUUID, logic.TotalAmount+logic.TricycleFee, logic.CreditAmount); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	// logic.SetFeilds()
	return
}

func (logic *OrderLogic) FillGoodsWeightWithTotal(db *gorm.DB) (err error) {
	list := logic.GoodsListRelated

	if len(list) == 0 {
		return nil
	}

	// 1. Collect valid goods UUIDs and filter out empty ones
	var goodsUUIDs []string
	validItems := make([]*model.BatchOrderGoods, 0, len(list))

	for _, item := range list {
		if item.GoodsUUID != "" {
			goodsUUIDs = append(goodsUUIDs, item.GoodsUUID)
			validItems = append(validItems, item)
		}
	}

	logic.GoodsListRelated = validItems

	if len(goodsUUIDs) == 0 {
		return nil
	}

	// 2. Use cache or database to get goods fields
	var goodsFieldMap map[string]model.GoodsFeild

	// Try cache first
	if goodsFieldMap, err = cache.GoodsCache.BatchGoodsFeildSet(goodsUUIDs, logic.OwnerUser); err != nil {
		return
	}
	// 3. Calculate and set fields for valid items only
	finalItems := make([]*model.BatchOrderGoods, 0, len(validItems))
	var totalAmount float64

	for _, item := range validItems {
		goodsField, ok := goodsFieldMap[item.GoodsUUID]
		if !ok {
			// Goods UUID doesn't exist, skip this item
			continue
		}

		// Set common fields
		item.BatchUUID = logic.BatchUUID
		item.OwnerUser = logic.OwnerUser
		item.UserUUID = logic.UserUUID
		item.GoodType = goodsField.GoodsTyp
		// item.GoodsTyp = goodsField.GoodsTyp

		var subTotal float64

		if goodsField.GoodsTyp == common.GoodsTypeFix {
			item.Weight = utils.FloatReserve(float64(item.Mount)*goodsField.GoodsWeight, 1)
			subTotal = utils.FloatReserve(float64(item.Mount)*item.Price, 0)
		}

		if goodsField.GoodsTyp == common.GoodsTypeBulk {
			item.Mount = 0
			subTotal = utils.FloatReserve(item.Price*item.Weight, 0)
		}

		item.GoodsFeild = goodsField
		item.Total = subTotal
		totalAmount += subTotal

		finalItems = append(finalItems, item)
	}

	logic.GoodsListRelated = finalItems
	logic.TotalAmount = totalAmount + logic.TricycleFee
	logic.CreditAmount = logic.TotalAmount - logic.Discount

	return nil
}

// 下单
func (logic *OrderLogic) Create(orderReq request.OrderReq) (err error) {
	// 假如是现金客户

	logic.BatchOrder = orderReq.BatchOrder

	// 前端未传下单时间时用当前时间
	if logic.OrderTime == 0 {
		logic.OrderTime = int(time.Now().Unix())
	}

	// 计算total
	if err = logic.FillGoodsWeightWithTotal(global.Global.DB); err != nil {
		return
	}

	logic.CreditAmount = logic.TotalAmount - logic.Discount - orderReq.FPayAmount
	payableAmount := logic.TotalAmount - logic.Discount
	if utils.FloatEqual(payableAmount, orderReq.FPayAmount) || utils.FloatGreat(orderReq.FPayAmount, payableAmount) {
		logic.Status = common.BatchOrderFinish
	} else {
		logic.Status = common.BatchOrderTemp
	}

	// 现金客户付款金额不能小于订单金额（赊欠单说明没付够）
	isCash, err := logic.isCashCustomer()
	if err != nil {
		return
	}
	if isCash && logic.Status == common.BatchOrderTemp {
		return common.CashCustomerPayLessErr
	}

	if logic.BatchUUID == "" {
		return common.BatchUUIDRequireErr
	}
	if len(logic.GoodsListRelated) == 0 {
		return common.BatchOrderGoodsRequireErr
	}

	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
	}

	// 用 Redis 幂等键防止网络抖动导致同一笔下单重复创建
	if !logic.acquireOrderDedup() {
		return common.OrderDuplicateSubmitErr
	}

	tx := logic.runtime.DB.Begin()

	// 1. 先创建订单（不包含货品）
	goodsList := logic.GoodsListRelated
	logic.GoodsListRelated = nil // 暂时清空

	if err = dao.OrderDao.Create(tx, &logic.BatchOrder); err != nil {
		logic.GoodsListRelated = goodsList // 恢复
		tx.Rollback()
		return
	}

	// 2. 再创建订单货品
	logic.GoodsListRelated = goodsList // 恢复
	for _, goods := range logic.GoodsListRelated {
		goods.BatchOrderUID = logic.UID
	}

	if err = dao.OrderDao.CreateGoodsBatch(tx, logic.OwnerUser, logic.GoodsListRelated); err != nil {
		tx.Rollback()
		return
	}

	// orderPay := NewOrderPayLogic(logic.context)
	// orderPay.BatchOrderUUID = logic.UID
	// // orderPay.Amount = order.TotalAmount - order.CreditAmount
	// orderPay.Amount = orderReq.FPayAmount
	// if err = orderPay.Create(tx, false); err != nil {
	// 	tx.Rollback()
	// 	return
	// }

	// 更新客户订单统计
	if err = dao.CustomerDao.IncrOrderStats(tx, logic.UserUUID, logic.TotalAmount+logic.TricycleFee, logic.CreditAmount); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	// if utils.FloatGreat(0.0, logic.CreditAmount) {
	// 	go logic.Record(false, model.HistoryStepCash, model.PayFeild{
	// 		PayFee:  utils.Violent2String(orderReq.TotalAmount),
	// 		PayType: orderReq.PayType,
	// 		PaidFee: utils.Violent2String(orderReq.FPayAmount),
	// 	})
	// } else {
	// 	go logic.Record(false, model.HistoryStepCredit, model.PayFeild{
	// 		PayFee:  utils.Violent2String(orderReq.TotalAmount),
	// 		PayType: orderReq.PayType,
	// 		PaidFee: utils.Violent2String(orderReq.FPayAmount),
	// 	})
	// }
	// logic.SetFeilds()
	return
}

func (logic *OrderLogic) Shared() (err error) {
	// 分享以客户为维度：同一客户当天分享多个订单只记录一条，更新最近一次分享时间
	if logic.UserUUID != "" {
		if err = dao.ShareDao.UpsertShare(logic.runtime.DB, &model.ShareRecord{
			OwnerUser:    logic.OwnerUser,
			CustomerUUID: logic.UserUUID,
			ShareTime:    time.Now().Unix(),
		}); err != nil {
			return
		}
	}
	// logic.Record(true, model.HistoryStepOrderShare, model.PayFeild{})
	return
}

// ShareDailyOrder builds the share data for today's orders of a customer.
// If orderUUID is provided, queries that order's date data.
// Otherwise, queries today's data.
func (logic *OrderLogic) ShareOrder(req request.ShareDailyOrderReq) (rsp *response.ShareDailyOrderRsp, err error) {
	// Get customer info for name
	customer, err := dao.CustomerDao.FromUUID(logic.runtime.DB, req.CustomerUUID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// 经营户（当前登录用户）的菜市场信息，用于分享单抬头与底部联系方式
	owner, err := dao.UserDao.FromUUID(logic.runtime.DB, logic.OwnerUser)
	if err != nil {
		return nil, fmt.Errorf("owner user not found: %w", err)
	}

	order, _, err := dao.OrderDao.GetShareOrdersWithGoods(logic.runtime.DB, logic.OwnerUser, req.CustomerUUID, req.OrderUUID)
	if err != nil {
		return
	}

	// Build goods name map via cache
	goodsUUIDs := make([]string, 0, len(order.GoodsListRelated))
	for _, g := range order.GoodsListRelated {
		goodsUUIDs = append(goodsUUIDs, g.GoodsUUID)
	}
	goodsMap, _ := cache.GoodsCache.BatchGoodsFeildSet(goodsUUIDs, logic.OwnerUser)

	// 只取最近的一单（同批次、不同批次都只展示最新的那个）
	var latestCredit float64
	var previousCreditAmount float64
	pendingDebt := customer.Debt - customer.PaidDebt
	if pendingDebt < 0 {
		pendingDebt = 0
	}

	// 表示查到了
	if order.UID != "" {
		latestCredit = order.CreditAmount

		preTime := order.CreatedAt
		if order.OrderTime > 0 {
			preTime = time.Unix(int64(order.OrderTime), 0)
		}

		// 历史单：精确回溯到该订单下单之前（只含订单赊欠，不含原始债务）
		previousCreditAmount, err = dao.OrderDao.CustomerCreditAmountBeforeOrder(
			logic.runtime.DB, logic.OwnerUser, req.CustomerUUID, preTime.Unix())
		if err != nil {
			global.Global.Logger.Error("[ShareOrder] 查询订单前欠失败，回退为客户总赊欠减本单赊欠", zap.Error(err))
			previousCreditAmount = customer.Credit - latestCredit
		} else {
			previousCreditAmount += pendingDebt
		}
	} else {
		previousCreditAmount += customer.Credit
	}

	rsp = &response.ShareDailyOrderRsp{
		PreviousCreditAmount: previousCreditAmount,
		TodayCreditAmount:    latestCredit,
	}

	// 渲染这一单的货品明细
	if order.UID != "" {
		ot := time.Unix(int64(order.OrderTime), 0)
		group := response.ShareDailyOrderGroup{
			OrderUUID:   order.UID,
			SerialNo:    utils.UUIDToSerial(order.UID),
			DateStr:     ot.Format("2006-01-02 15:04"),
			TricycleFee: order.TricycleFee,
			TotalAmount: order.CreditAmount,
		}
		for _, g := range order.GoodsListRelated {
			item := &response.ShareDailyOrderItem{
				GoodsUUID: g.GoodsUUID,
				GoodsName: goodsMap[g.GoodsUUID].GoodsName,
				Price:     g.Price,
				Weight:    g.Weight,
				Mount:     g.Mount,
				Total:     g.Total,
				GoodType:  g.GoodType,
			}
			group.GoodsList = append(group.GoodsList, item)
			group.TotalMount += g.Mount
			if g.GoodType == common.GoodsTypeBulk {
				group.TotalWeight += g.Weight
			}
		}
		// 先散装后定装（GoodType 降序：散装=2 在前，定装=1 在后），
		// 同类型内按货品名升序，保证不同时段分享次序稳定
		sort.SliceStable(group.GoodsList, func(i, j int) bool {
			if group.GoodsList[i].GoodType != group.GoodsList[j].GoodType {
				return group.GoodsList[i].GoodType > group.GoodsList[j].GoodType
			}
			return group.GoodsList[i].GoodsName < group.GoodsList[j].GoodsName
		})
		rsp.OrderGroups = append(rsp.OrderGroups, group)

		rsp.Date = ot.Format("2006-01-02 15:04")
	} else {
		rsp.Date = time.Now().Format("2006-01-02 15:04")
	}

	// 计算总赊欠 = 前欠 + 今日赊欠
	rsp.TotalCreditAmount = rsp.PreviousCreditAmount + rsp.TodayCreditAmount

	// 抬头：<省份>食用农产品批发市场销售票据 / 市场名称 / 用户名 / 档口名称
	province := strings.TrimSpace(owner.Province)
	rsp.Title = province + "食用农产品批发市场销售票据"
	rsp.MarketName = owner.MarketName
	rsp.UserName = owner.Name
	rsp.StallName = owner.StallName
	// 底部联系方式：优先用户自定义的分享联系方式（前端逗号分隔，展示用空格分隔），
	// 未设置时回退到登录手机号，最后回退到配置文件
	rsp.Contact = joinContactPhones(owner.ContactPhones)
	if rsp.Contact == "" {
		rsp.Contact = owner.Phone
	}
	if rsp.Contact == "" {
		rsp.Contact = logic.runtime.Viper.GetString("share.contact")
	}
	rsp.CustomerName = customer.Name
	return
}

// joinContactPhones 将前端用逗号分隔的多个联系方式转成展示用的空格分隔
func joinContactPhones(raw string) string {
	parts := strings.Split(raw, ",")
	phones := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			phones = append(phones, part)
		}
	}
	return strings.Join(phones, " ")
}

// buildReceiptData 从 ShareDailyOrderRsp 构建模板渲染数据，按订单分组。
func (logic *OrderLogic) buildReceiptData(rsp *response.ShareDailyOrderRsp) utils.ReceiptData {
	var orders []utils.ReceiptOrder
	for _, og := range rsp.OrderGroups {
		var goodsItems []utils.ReceiptGoodsItem
		for _, g := range og.GoodsList {
			goodsItems = append(goodsItems, utils.ReceiptGoodsItem{
				GoodsName: g.GoodsName,
				MountStr:  utils.FormatMount(g.Mount),
				WeightStr: utils.FormatWeight(g.Weight),
				PriceStr:  utils.FormatFloat(g.Price),
				TotalStr:  utils.FormatFloat(g.Total),
				GoodType:  g.GoodType,
			})
		}
		orders = append(orders, utils.ReceiptOrder{
			SerialNo:       og.SerialNo,
			Date:           og.DateStr,
			GoodsList:      goodsItems,
			TotalMount:     fmt.Sprintf("%d", og.TotalMount),
			TotalWeight:    utils.FormatTotalWeight(og.TotalWeight),
			TotalAmountStr: utils.FormatFloat(og.TotalAmount),
			TricycleFeeStr: utils.FormatFloat(og.TricycleFee),
		})
	}
	return utils.ReceiptData{
		Title:                  rsp.Title,
		CustomerName:           rsp.CustomerName,
		SerialNo:               fmt.Sprintf("%d", time.Now().Unix()%100000),
		Date:                   rsp.Date,
		Contact:                rsp.Contact,
		MarketName:             rsp.MarketName,
		UserName:               rsp.UserName,
		StallName:              rsp.StallName,
		Orders:                 orders,
		CreditAmountStr:        utils.FormatFloat(rsp.TodayCreditAmount),
		TotalPreviousCreditStr: utils.FormatFloat(rsp.PreviousCreditAmount),
		TotalDebtStr:           utils.FormatFloat(rsp.TotalCreditAmount),
	}
}

// GenerateShareHTML builds the receipt HTML for a customer's daily order and returns the HTML string.
func (logic *OrderLogic) GenerateShareHTML(req request.ShareDailyOrderReq) (html string, err error) {
	rsp, err := logic.ShareOrder(req)
	if err != nil {
		return
	}

	receiptData := logic.buildReceiptData(rsp)
	html, err = utils.GenerateHTML(receiptData)
	if err != nil {
		logic.runtime.Logger.Error("GenerateShareHTML failed", zap.Error(err))
	}
	return
}

func (logic *OrderLogic) UpdateStatus() (err error) {
	// 先加载当前订单数据，用于调整客户统计
	var oldOrder model.BatchOrder
	if oldOrder, err = logic.LoadSingle(logic.UID); err != nil {
		return
	}

	// 订单状态更新 + 客户统计回退放在同一个事务内
	tx := logic.runtime.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = dao.OrderDao.UpdateStatus(tx, logic.OwnerUser, logic.UID, logic.Status); err != nil {
		return
	}

	// 作废/退款/退货：减少客户统计
	if logic.Status == common.BatchOrderCancel || logic.Status == common.BatchOrderRefund || logic.Status == common.BatchOrderReTurn {
		if err = dao.CustomerDao.DecrOrderStats(tx, oldOrder.UserUUID, oldOrder.TotalAmount, oldOrder.CreditAmount); err != nil {
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		return
	}
	return
}

// 更新单次
func (logic *OrderLogic) Update(orderReq request.OrderReq) (err error) {
	var old model.BatchOrder
	if old, err = logic.LoadSingle(orderReq.UID); err != nil {
		return
	}
	logic.BatchOrder = orderReq.BatchOrder // owner_user重置成了空字符串
	logic.BatchUUID = old.BatchUUID

	if logic.BatchOrder.OrderTime == 0 {
		logic.BatchOrder.OrderTime = old.OrderTime
	}

	if logic.BatchOrder.UserUUID == "" {
		logic.BatchOrder.UserUUID = old.UserUUID
	}

	// 计算total
	if err = logic.FillGoodsWeightWithTotal(global.Global.DB); err != nil {
		return
	}

	logic.CreditAmount = logic.TotalAmount - logic.Discount - orderReq.FPayAmount
	payableAmount := logic.TotalAmount - logic.Discount
	if utils.FloatEqual(payableAmount, orderReq.FPayAmount) || utils.FloatGreat(orderReq.FPayAmount, payableAmount) {
		logic.Status = common.BatchOrderFinish
	} else {
		logic.Status = common.BatchOrderTemp
	}

	// 现金客户付款金额不能小于订单金额（赊欠单说明没付够）
	isCash, err := logic.isCashCustomer()
	if err != nil {
		return
	}
	if isCash && logic.Status == common.BatchOrderTemp {
		return common.CashCustomerPayLessErr
	}

	tx := logic.runtime.DB.Begin()
	// 分表删除：删除原货品记录（走分表，避免外键约束）
	if err = dao.OrderDao.DeleteGoodsByOrderUUID(tx, logic.OwnerUser, logic.UID); err != nil {
		tx.Rollback()
		return err
	}
	// 分表删除：删除原还款记录
	if err = dao.OrderDao.DeletePayByOrderUUID(tx, logic.OwnerUser, logic.UID); err != nil {
		tx.Rollback()
		return err
	}
	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
		goods.BatchOrderUID = logic.UID
	}
	// 分表更新订单（使用 Updates 而非 Save，避免 GORM 处理 GoodsListRelated 关联触犯主表外键约束）
	// 注意：OrderLogic 同时定义了 Shared() 方法和嵌入的 BatchOrder.Shared 字段，
	// 直接用 logic.Shared 会被解析成方法值 func() error，必须显式访问字段。
	updateFields := map[string]interface{}{
		"updated_at": time.Now(),
		// "share_time":    logic.BatchOrder.SharedTime,
		"status":        logic.BatchOrder.Status,
		"amount":        logic.BatchOrder.TotalAmount,
		"credit_amount": logic.BatchOrder.CreditAmount,
		"tricycle_fee":  logic.BatchOrder.TricycleFee,
		"order_time":    logic.BatchOrder.OrderTime,
		"user_uuid":     logic.BatchOrder.UserUUID,
	}
	if err = dao.OrderDao.UpdateOrderFields(tx, logic.OwnerUser, logic.UID, updateFields); err != nil {
		tx.Rollback()
		return
	}

	// 分表批量创建新货品记录
	if err = dao.OrderDao.CreateGoodsBatch(tx, logic.OwnerUser, logic.GoodsListRelated); err != nil {
		tx.Rollback()
		return
	}

	orderPay := NewOrderPayLogic(logic.context)
	orderPay.BatchOrderUUID = logic.UID
	// orderPay.Amount = order.TotalAmount - order.CreditAmount
	orderPay.Amount = orderReq.FPayAmount
	if err = orderPay.Create(tx, false); err != nil {
		return
	}

	// 更新客户订单统计
	if logic.UserUUID == old.UserUUID {
		// 同客户：按差额调整（新金额 - 旧金额）
		deltaTotal := logic.TotalAmount - old.TotalAmount
		deltaDebt := logic.CreditAmount - old.CreditAmount
		if deltaTotal != 0 || deltaDebt != 0 {
			if err = dao.CustomerDao.UpdateOrderStatsDiff(tx, logic.UserUUID, deltaTotal, deltaDebt); err != nil {
				tx.Rollback()
				return
			}
		}
	} else {
		// 换客户：旧客户回退统计（订单数-1、总金额-旧、credit-旧赊欠），新客户累加统计
		if err = dao.CustomerDao.DecrOrderStats(tx, old.UserUUID, old.TotalAmount, old.CreditAmount); err != nil {
			tx.Rollback()
			return
		}
		if err = dao.CustomerDao.IncrOrderStats(tx, logic.UserUUID, logic.TotalAmount, logic.CreditAmount); err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	// if utils.FloatGreat(0.0, logic.CreditAmount) {
	// 	go logic.Record(false, model.HistoryStepCash, model.PayFeild{
	// 		PayFee:  utils.Violent2String(orderReq.TotalAmount),
	// 		PayType: orderReq.PayType,
	// 		PaidFee: utils.Violent2String(orderReq.FPayAmount),
	// 	})
	// } else {
	// 	go logic.Record(false, model.HistoryStepCredit, model.PayFeild{
	// 		PayFee:  utils.Violent2String(orderReq.TotalAmount),
	// 		PayType: orderReq.PayType,
	// 		PaidFee: utils.Violent2String(orderReq.FPayAmount),
	// 	})
	// }

	// logic.SetFeilds()
	// logic.Record(false, model.HistoryStepOrderFix, model.PayFeild{})
	return
}

func (logic *OrderLogic) FromUUID(uuid string) (err error) {
	if logic.BatchOrder, err = dao.OrderDao.FindByUID(logic.runtime.DB, logic.OwnerUser, uuid, "GoodsListRelated"); err != nil {
		return
	}
	logic.LoadHistory()
	logic.SetFeilds()
	// logic.SetTotal()
	return
}

func (logic *OrderLogic) LoadSingle(uuid string) (batchOrder model.BatchOrder, err error) {
	batchOrder, err = dao.OrderDao.FindByUID(logic.runtime.DB, logic.OwnerUser, uuid)
	return
}

func (logic *OrderLogic) FindLatestGoods(goodsUUIDList []string) (goodsOrderList []response.BatchOrderGoodsOrderRsp, err error) {
	for _, goodsUUID := range goodsUUIDList {
		conds := []utils.Cond{
			utils.NewWhereCond("goods_uuid", goodsUUID),
			utils.CreatedOrderAscCond(),
		}
		// todo BatchOrderGoods需要加索引
		var goodsOrder model.BatchOrderGoods
		var _price float64
		if err = utils.First(logic.runtime.DB.Table(dao.ShardOrderGoodsTable(logic.OwnerUser)), &goodsOrder, conds...); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var goods model.Goods
				if err = utils.First(logic.runtime.DB, &goods, utils.NewWhereCond("uid", goodsUUID)); err != nil {
					return
				}
				_price = goods.Price
			} else {
				return
			}
		} else {
			_price = goodsOrder.Price
		}
		goodsOrderList = append(goodsOrderList, response.BatchOrderGoodsOrderRsp{
			GoodsUUID: goodsUUID,
			Price:     _price,
		})
	}

	return
}

func (logic *OrderLogic) SetGoodsFeild() (err error) {
	var goodsUUIDList []string
	for _, goods := range logic.GoodsListRelated {
		goodsUUIDList = append(goodsUUIDList, goods.GoodsUUID)
	}

	goodsM := make(map[string]model.GoodsFeild, 0)
	if goodsM, err = cache.GoodsCache.BatchGoodsFeildSet(goodsUUIDList, logic.OwnerUser); err != nil {
		return
	}
	for _, goods := range logic.GoodsListRelated {
		goods.GoodsFeild = goodsM[goods.GoodsUUID]
	}

	return
}

func (logic *OrderLogic) SetCustomerFeild() (err error) {
	logic.CustomerFeild, err = cache.CustomerCache.CustomerFeildSet(logic.UserUUID, logic.OwnerUser)
	return
}

func (logic *OrderLogic) List(req request.BatchOrderListReq) (orderList []*model.BatchOrder, err error) {
	// 使用分表 DAO 查询
	var startT, endT time.Time
	if req.StartTime != 0 {
		startT = time.Unix(req.StartTime, 0)
	}
	if req.EndTime != 0 {
		endT = time.Unix(req.EndTime, 0)
	}

	limitCond := utils.DefaultSetLimitCond(req.LimitCond)
	offset := 0
	if limitCond.Page > 0 && limitCond.PageCount > 0 {
		offset = (limitCond.Page - 1) * limitCond.PageCount
	}

	// 收集需要查询的客户 UUID 列表
	var userUUIDs []string
	if req.UserNameSearch != "" {
		// 先根据用户名模糊查询客户 UUID 列表
		userUUIDs, err = dao.CustomerDao.FindUserUUIDsByName(logic.runtime.DB, logic.OwnerUser, req.UserNameSearch)
		if err != nil {
			return
		}
		if len(userUUIDs) == 0 {
			// 没有匹配的客户，直接返回空列表
			return
		}
	} else if req.UserUUID != "" {
		userUUIDs = []string{req.UserUUID}
	}

	orderList, err = dao.OrderDao.List(
		logic.runtime.DB,
		logic.OwnerUser,
		userUUIDs,
		startT,
		endT,
		req.Status,
		offset,
		limitCond.PageCount,
		true, // 预加载货品
	)

	if err != nil {
		return
	}

	logic.BatchFeilds(orderList)

	return
}

func (logic *OrderLogic) BatchFeilds(orderList []*model.BatchOrder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var _cl []string
		for _, o := range orderList {
			_cl = append(_cl, o.UserUUID)
		}
		_cm, _ := cache.CustomerCache.BatchCustomerFeildSet(_cl, logic.OwnerUser)

		for _, o := range orderList {
			o.CustomerFeild = _cm[o.UserUUID]
		}
		wg.Done()
	}()
	go func() {
		var _gl []string
		for _, o := range orderList {
			for _, g := range o.GoodsListRelated {
				_gl = append(_gl, g.GoodsUUID)
			}
		}
		_gm, _ := cache.GoodsCache.BatchGoodsFeildSet(_gl, logic.OwnerUser)
		for _, o := range orderList {
			for _, g := range o.GoodsListRelated {
				g.GoodsFeild = _gm[g.GoodsUUID]
				g.SetTotal()
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

// func (logic *OrderLogic) SetTotal() {
// 	for _, g := range logic.GoodsListRelated {
// 		logic.TotalAmount += g.SetTotal()
// 	}
// }

func (logic *OrderLogic) BatchSetFeilds(orders []model.BatchOrder) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		logic.SetCustomerFeild()
		wg.Done()
	}()
	go func() {
		logic.SetGoodsFeild()
		wg.Done()
	}()
	wg.Wait()
}

// func (logic *OrderLogic) BatchSetFeilds(orders []model.BatchOrder) {
// 	wg := sync.WaitGroup{}
// 	wg.Add(2)
// 	go func() {
// 		logic.SetCustomerFeild()
// 		wg.Done()
// 	}()
// 	go func() {
// 		logic.SetGoodsFeild()
// 		wg.Done()
// 	}()
// 	wg.Wait()
// }

func (logic *OrderLogic) SetFeilds() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		logic.SetCustomerFeild()
		wg.Done()
	}()
	go func() {
		logic.SetGoodsFeild()
		wg.Done()
	}()
	wg.Wait()
}

func (logic *OrderLogic) Record(loadself bool, stepType int32, pay model.PayFeild) {
	go func() {
		if loadself {
			if err := utils.Find(logic.runtime.DB.Preload("GoodsListRelated"), &logic.BatchOrder); err != nil {
				logic.runtime.Logger.Error(fmt.Sprintf("OrderLogic Record: %s", err))
				return
			}
			logic.SetFeilds()
		}
		logic.BatchOrder.Record(logic.runtime.DB, stepType, pay)
	}()
}

func (logic *OrderLogic) LoadHistory() (err error) {
	var history model.BatchOrderHistory
	if err = utils.First(logic.runtime.DB, &history, utils.NewWhereCond("batch_order_uuid", logic.UID)); err != nil {
		logic.runtime.Logger.Error("OrderLogic LoadHistory", zap.Error(err))
		return
	}
	if history.UID != "" {
		logic.History = history
	}
	return
}

// func (logic *OrderLogic) GoodsList(req request.BatchGoodsListReq) (rsp response.BatchGoodsGroupRsp, err error) {
// 	rspItems := make([]*response.BatchGoodsGroupItem, 0)
// 	var storages []model.BatchOrderGoods
// 	db := logic.runtime.DB

// 	conds := []model.Cond{
// 		model.NewWhereCond("batch_uuid", req.BatchUUID),
// 		model.NewWhereCond("goods_uuid", req.GoodsUUID),
// 	}

// 	conds = append(conds, model.CreatedOrderDescCond())
// 	if err = model.Find(db, &storages, conds...); err != nil {
// 		logic.runtime.Logger.Error("BatchLogic List FindBatch", zap.Any("req", req), zap.Error(err))
// 		return
// 	}

// 	var customerUUIDList []string
// 	for _, storage := range storages {
// 		customerUUIDList = append(customerUUIDList, storage.UserUUID)
// 	}

// 	_cm, _ := model.BatchCustomerFeildSet(logic.runtime.DB, customerUUIDList, logic.OwnerUser)

//		for _, storage := range storages {
//			rspItems = append(rspItems, &response.BatchGoodsGroupItem{
//				SellTime:     storage.CreatedAt.Format("2006-01-02 15:04:05"),
//				SellPrice:    utils.FloatReserveStr(storage.Price, 1), //
//				SellAmount:   utils.FloatReserveStr(storage.Amount(), 1),
//				CustomerName: _cm[storage.UserUUID].CustomerName,
//				SellTotall:   utils.FloatReserveStr(float64(storage.Price)*storage.Amount(), 1),
//				GoodType:     storage.GoodType,
//			})
//		}
//		rsp.Items = rspItems
//		return
//	}
func (logic *OrderLogic) GoodsList(req request.BatchGoodsListReq) (rsp response.BatchGoodsGroupRsp, err error) {
	rspItems := make([]*response.BatchGoodsGroupItem, 0)

	storages, err := dao.OrderDao.FindGoodsByBatchAndGoodsUUID(logic.runtime.DB, logic.OwnerUser, req.BatchUUID, req.GoodsUUID, common.ValidOrder)
	if err != nil {
		logic.runtime.Logger.Error("GoodsList FindGoodsByBatchAndGoodsUUID", zap.Any("req", req), zap.Error(err))
		return
	}

	if len(storages) == 0 {
		return
	}

	// 1. 收集客户UUID并获取客户信息
	var customerUUIDList []string
	for _, storage := range storages {
		customerUUIDList = append(customerUUIDList, storage.UserUUID)
	}
	_cm, _ := cache.CustomerCache.BatchCustomerFeildSet(customerUUIDList, logic.OwnerUser)

	// 2. 查询对应的BatchGoods获取成本信息
	batchGoods, err := dao.BatchGoodsDao.FromUUID(logic.runtime.DB, logic.OwnerUser, req.BatchUUID, req.GoodsUUID)
	if err != nil {
		logic.runtime.Logger.Error("getBatchGoods failed", zap.Error(err))
		return
	}

	// 3. 计算总销售额、总成本、总销售数量/重量
	var totalSellAmount float64   // 总销售额
	var totalCostAmount float64   // 总成本
	var totalSellQuantity float64 // 总销售数量/重量
	var totalSellMount int        // 定装总件数

	for _, storage := range storages {
		// 获取客户名称
		customerName := ""
		if customer, ok := _cm[storage.UserUUID]; ok {
			customerName = customer.CustomerName
		}

		// 计算单项销售额
		sellAmount := storage.Sell()
		sellTotal := utils.FloatReserve(float64(storage.Price)*sellAmount, 1)

		// 累加总销售额
		totalSellAmount += sellTotal

		// 根据类型累加销售数量和成本
		if storage.GoodType == common.GoodsTypeFix {
			// 定装：按件数计算
			totalSellMount += storage.Mount
			totalCostAmount += float64(storage.Mount) * batchGoods.Price
		} else if storage.GoodType == common.GoodsTypeBulk {
			totalSellQuantity += storage.Weight
			totalCostAmount += storage.Weight * batchGoods.Price
		}

		rspItems = append(rspItems, &response.BatchGoodsGroupItem{
			SellTime:     storage.CreatedAt.Format("2006-01-02 15:04:05"),
			SellPrice:    utils.FloatReserveStr(storage.Price, 1),
			SellAmount:   utils.FloatReserveStr(sellAmount, 1),
			CustomerName: customerName,
			SellTotall:   utils.FloatReserveStr(sellTotal, 1),
			GoodType:     storage.GoodType,
		})
	}

	// 4. 计算利润
	profit := totalSellAmount - totalCostAmount
	rsp.Profit = utils.FloatReserveStr(profit, 1)

	// 5. 计算剩余库存
	var surplus float64
	if batchGoods != nil {
		if batchGoods.GoodType == common.GoodsTypeFix {
			// 定装：原始件数 - 已卖件数
			surplus = float64(batchGoods.Mount - totalSellMount)
		} else if batchGoods.GoodType == common.GoodsTypeBulk {
			// 散装：原始重量 - 已卖重量
			surplus = batchGoods.Weight - totalSellQuantity
		}
	}
	if batchGoods.GoodType == common.GoodsTypeFix {
		rsp.Surplus = utils.FloatReserveStr(surplus, 1) + "件"
	} else if batchGoods.GoodType == common.GoodsTypeBulk {
		rsp.Surplus = utils.FloatReserveStr(surplus, 1) + "斤"
	}

	rsp.Items = rspItems
	return
}

// CreditListLogic 赊欠列表业务逻辑
func (logic *OrderLogic) CreditList(listReq request.CreditListReq) (rsp *response.CreditListResponse, err error) {
	if rsp, err = dao.OrderDao.GetCreditList(logic.runtime.DB, logic.OwnerUser, listReq); err != nil {
		return
	}

	if len(rsp.List) == 0 {
		return
	}

	var userUUIDList []string

	for _, item := range rsp.List {
		userUUIDList = append(userUUIDList, item.UserUUID)
	}

	userMapper, _ := cache.CustomerCache.BatchCustomerFeildSet(userUUIDList, logic.OwnerUser)

	for _, item := range rsp.List {
		item.UserName = userMapper[item.UserUUID].CustomerName
	}

	return
}
