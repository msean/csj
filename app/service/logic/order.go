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
	"errors"
	"fmt"
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

// 码单
func (logic *OrderLogic) TempCreate() (err error) {
	logic.BatchOrder.DefaultSet()

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

	// 创建订单 + 更新客户统计（使用事务保证一致性）
	tx := logic.runtime.DB.Begin()
	if err = utils.CreateObj(tx, &logic.BatchOrder); err != nil {
		tx.Rollback()
		global.Global.Logger.Error("OrderLogic TempCreate CreateObj", zap.Any("logic.BatchOrder", logic.BatchOrder), zap.Error(err))
		return
	}

	if err = dao.CustomerDao.IncrOrderStats(tx, logic.UserUUID, logic.TotalAmount); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	logic.SetFeilds()
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
	fmt.Println(">>>>>>>>>>OwnerUser2", logic.OwnerUser)
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
	logic.TotalAmount = totalAmount

	return nil
}

// 下单
func (logic *OrderLogic) Create(orderReq request.OrderReq) (err error) {
	logic.BatchOrder = orderReq.BatchOrder
	// 计算total
	if err = logic.FillGoodsWeightWithTotal(global.Global.DB); err != nil {
		return
	}

	logic.CreditAmount = logic.TotalAmount - orderReq.FPayAmount
	if utils.FloatEqual(logic.TotalAmount, orderReq.FPayAmount) || utils.FloatGreat(orderReq.FPayAmount, logic.TotalAmount) {
		logic.Status = common.BatchOrderFinish
	} else {
		logic.Status = common.BatchOrderTemp
	}

	logic.BatchOrder.Shared = common.BatchOrderUnshare

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

	tx := logic.runtime.DB.Begin()
	if err = utils.CreateObj(tx, &logic.BatchOrder); err != nil {
		tx.Rollback()
		return
	}

	orderPay := NewOrderPayLogic(logic.context)
	orderPay.BatchOrderUUID = logic.UID
	// orderPay.Amount = order.TotalAmount - order.CreditAmount
	orderPay.Amount = orderReq.FPayAmount
	if err = orderPay.Create(tx, false); err != nil {
		tx.Rollback()
		return
	}

	// 更新客户订单统计
	if err = dao.CustomerDao.IncrOrderStats(tx, logic.UserUUID, logic.TotalAmount); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	if utils.FloatGreat(0.0, logic.CreditAmount) {
		go logic.Record(false, model.HistoryStepCash, model.PayFeild{
			PayFee:  utils.Violent2String(orderReq.TotalAmount),
			PayType: orderReq.PayType,
			PaidFee: utils.Violent2String(orderReq.FPayAmount),
		})
	} else {
		go logic.Record(false, model.HistoryStepCredit, model.PayFeild{
			PayFee:  utils.Violent2String(orderReq.TotalAmount),
			PayType: orderReq.PayType,
			PaidFee: utils.Violent2String(orderReq.FPayAmount),
		})
	}
	logic.SetFeilds()
	return
}

func (logic *OrderLogic) Shared() (err error) {
	if err = dao.OrderDao.Shared(logic.runtime.DB, logic.UID); err != nil {
		return
	}
	logic.Record(true, model.HistoryStepOrderShare, model.PayFeild{})
	return
}

func (logic *OrderLogic) UpdateStatus() (err error) {
	// 先加载当前订单数据，用于调整客户统计
	var oldOrder model.BatchOrder
	if oldOrder, err = logic.LoadSingle(logic.UID); err != nil {
		return
	}

	if err = dao.OrderDao.UpdateStatus(logic.runtime.DB, logic.UID, logic.Status); err != nil {
		return
	}

	// 作废/退款/退货：减少客户统计
	if logic.Status == common.BatchOrderCancel || logic.Status == common.BatchOrderRefund || logic.Status == common.BatchOrderReTurn {
		if err = dao.CustomerDao.DecrOrderStats(logic.runtime.DB, logic.UserUUID, oldOrder.TotalAmount); err != nil {
			return
		}
	}

	switch logic.Status {
	case common.BatchOrderTemp:
		logic.Record(true, model.HistoryStepCredit, model.PayFeild{})
	case common.BatchOrderCancel, common.BatchOrderRefund:
		logic.Record(true, model.HistoryStepCrash, model.PayFeild{})
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

	// 计算total
	if err = logic.FillGoodsWeightWithTotal(global.Global.DB); err != nil {
		return
	}

	logic.CreditAmount = logic.TotalAmount - orderReq.FPayAmount
	if utils.FloatEqual(logic.TotalAmount, orderReq.FPayAmount) || utils.FloatGreat(orderReq.FPayAmount, logic.TotalAmount) {
		logic.Status = common.BatchOrderFinish
	} else {
		logic.Status = common.BatchOrderTemp
	}

	tx := logic.runtime.DB.Begin()
	if err = tx.Delete(&model.BatchOrderGoods{}, "batch_order_uuid=?", logic.UID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Delete(&model.BatchOrderPay{}, "batch_order_uuid=?", logic.UID).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, goods := range logic.GoodsListRelated {
		goods.OwnerUser = logic.OwnerUser
		goods.BatchUUID = logic.BatchUUID
		goods.UserUUID = logic.UserUUID
		goods.BatchOrderUID = logic.UID
	}
	if err = tx.Omit("created_at").Save(&logic.BatchOrder).Error; err != nil {
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

	// 更新客户订单统计差额（新金额 - 旧金额）
	deltaTotal := logic.TotalAmount - old.TotalAmount
	if deltaTotal != 0 {
		if err = dao.CustomerDao.UpdateOrderStatsDiff(tx, logic.UserUUID, deltaTotal); err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	if utils.FloatGreat(0.0, logic.CreditAmount) {
		go logic.Record(false, model.HistoryStepCash, model.PayFeild{
			PayFee:  utils.Violent2String(orderReq.TotalAmount),
			PayType: orderReq.PayType,
			PaidFee: utils.Violent2String(orderReq.FPayAmount),
		})
	} else {
		go logic.Record(false, model.HistoryStepCredit, model.PayFeild{
			PayFee:  utils.Violent2String(orderReq.TotalAmount),
			PayType: orderReq.PayType,
			PaidFee: utils.Violent2String(orderReq.FPayAmount),
		})
	}

	// logic.SetFeilds()
	// logic.Record(false, model.HistoryStepOrderFix, model.PayFeild{})
	return
}

func (logic *OrderLogic) FromUUID(uuid string) (err error) {
	if err = utils.Find(logic.runtime.DB.Preload("GoodsListRelated"), &logic.BatchOrder, utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	logic.LoadHistory()
	logic.SetFeilds()
	// logic.SetTotal()
	return
}

func (logic *OrderLogic) LoadSingle(uuid string) (batchOrder model.BatchOrder, err error) {
	err = utils.Find(logic.runtime.DB, &batchOrder, utils.WhereUIDCond(uuid))
	// logic.SetTotal()
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
		if err = utils.First(logic.runtime.DB, &goodsOrder, conds...); err != nil {
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

func (logic *OrderLogic) List(userUUID string, startTime, endTime int64, status int32, limitCond utils.LimitCond) (orderList []*model.BatchOrder, err error) {
	conds := []utils.Cond{
		limitCond,
		utils.NewWhereCond("owner_user", logic.OwnerUser),
	}
	if userUUID != "" {
		conds = append(conds, utils.NewWhereCond("user_uuid", userUUID))
	}

	if startTime != 0 {
		conds = append(conds, utils.NewCmpCond("created_at", ">=", time.Unix(startTime, 0)))
	}
	if endTime != 0 {
		conds = append(conds, utils.NewCmpCond("created_at", "<=", time.Unix(endTime, 0)))
	}
	if status != 0 {
		conds = append(conds, utils.NewWhereCond("status", status))
	} else {
		// conds = append(conds, model.NewInCondFromInt("status", common.ExCludeTempBatchOrder))
	}
	conds = append(conds, utils.CreatedOrderDescCond())
	if err = utils.Find(logic.runtime.DB.Preload("GoodsListRelated"), &orderList, conds...); err != nil {
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
	var storages []model.BatchOrderGoods
	db := logic.runtime.DB

	conds := []utils.Cond{
		utils.NewWhereCond("batch_uuid", req.BatchUUID),
		utils.NewWhereCond("goods_uuid", req.GoodsUUID),
	}
	conds = append(conds, utils.CreatedOrderDescCond())

	if err = utils.Find(db, &storages, conds...); err != nil {
		logic.runtime.Logger.Error("BatchLogic List FindBatch", zap.Any("req", req), zap.Error(err))
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
	rsp.Profits = utils.FloatReserveStr(profit, 1)

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
