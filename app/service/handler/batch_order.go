package handler

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func batchOrderRouter(g *gin.RouterGroup) {
	batchOrderGroup := g.Group("/batch_order", middleware.AuthMiddleware())
	{
		// 码单
		batchOrderGroup.POST("/temp/create", BatchOrderTempCreate)
		// 下单
		batchOrderGroup.POST("/create", BatchOrderCreate)
		batchOrderGroup.POST("/update", BatchOrderUpdate)
		batchOrderGroup.POST("/update_status", BatchOrderUpdateStatus)
		batchOrderGroup.POST("/detail", BatchOrderDetail)
		batchOrderGroup.POST("/list", BatchOrderGoodsList)
		batchOrderGroup.POST("/share", BatchOrderShared)
		batchOrderGroup.POST("/latest_by_goods", BatchOrderGoodsLatest)
	}
}

// 码单
func BatchOrderTempCreate(c *gin.Context) {
	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := order.TempCreate(); err != nil {
		common.Response(c, err, nil)
		return
	}
	go order.Record(false, model.HistoryStepOrder, model.PayFeild{})
	common.Response(c, nil, order)
}

func BatchOrderCreate(c *gin.Context) {
	type Form struct {
		model.BatchOrder
		FPayAmount float64 `json:"payAmount"` // 付款金额
		// FCreditAmount float64 `json:"creditAmount"`
		PayType int32 `json:"payType"` // 支付方式
	}
	var form Form
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	order.BatchOrder = form.BatchOrder // owner_user重置成了空字符串
	order.OwnerUser = common.GetUserUUID(c)

	// 计算total
	if err := order.FillGoodsWeightWithTotal(global.Global.DB); err != nil {
		common.Response(c, err, nil)
		return
	}

	tx := global.Global.DB.Begin()
	order.CreditAmount = order.TotalAmount - form.FPayAmount
	if utils.FloatEqual(order.TotalAmount, form.FPayAmount) || utils.FloatGreat(form.FPayAmount, order.TotalAmount) {
		order.Status = common.BatchOrderFinish
	} else {
		order.Status = common.BatchOrderTemp
	}

	if err := order.Create(tx); err != nil {
		common.Response(c, err, nil)
		return
	}
	orderPay := logic.NewBatchOrderPayLogic(c)
	orderPay.BatchOrderUUID = order.UID
	// orderPay.Amount = order.TotalAmount - order.CreditAmount
	orderPay.Amount = form.FPayAmount
	if err := orderPay.Create(tx, false); err != nil {
		common.Response(c, err, nil)
		return
	}
	tx.Commit()
	if utils.FloatGreat(0.0, order.CreditAmount) {
		go order.Record(false, model.HistoryStepCash, model.PayFeild{
			PayFee:  utils.Violent2String(form.TotalAmount),
			PayType: form.PayType,
			PaidFee: utils.Violent2String(form.FPayAmount),
		})
	} else {
		go order.Record(false, model.HistoryStepCredit, model.PayFeild{
			PayFee:  utils.Violent2String(form.TotalAmount),
			PayType: form.PayType,
			PaidFee: utils.Violent2String(form.FPayAmount),
		})
	}

	common.Response(c, nil, order)
}

func BatchOrderUpdate(c *gin.Context) {
	type Form struct {
		model.BatchOrder
		FPayAmount float64 `json:"payAmount"` // 付款金额
		PayType    int32   `json:"payType"`   // 支付方式
	}
	var err error
	var form Form
	if err = c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	var old model.BatchOrder
	if old, err = order.LoadSingle(form.UID); err != nil {
		return
	}
	order.BatchOrder = form.BatchOrder // owner_user重置成了空字符串
	order.BatchUUID = old.BatchUUID
	order.OwnerUser = common.GetUserUUID(c)

	// 计算total
	if err := order.FillGoodsWeightWithTotal(global.Global.DB); err != nil {
		common.Response(c, err, nil)
		return
	}

	tx := global.Global.DB.Begin()
	order.CreditAmount = order.TotalAmount - form.FPayAmount
	if utils.FloatEqual(order.TotalAmount, form.FPayAmount) || utils.FloatGreat(form.FPayAmount, order.TotalAmount) {
		order.Status = common.BatchOrderFinish
	} else {
		order.Status = common.BatchOrderTemp
	}

	// if err := order.Create(tx); err != nil {
	// 	common.Response(c, err, nil)
	// 	return
	// }

	if err := order.Update(tx); err != nil {
		common.Response(c, err, nil)
		return
	}
	orderPay := logic.NewBatchOrderPayLogic(c)
	orderPay.BatchOrderUUID = order.UID
	// orderPay.Amount = order.TotalAmount - order.CreditAmount
	orderPay.Amount = form.FPayAmount
	if err := orderPay.Create(tx, false); err != nil {
		common.Response(c, err, nil)
		return
	}
	tx.Commit()
	if utils.FloatGreat(0.0, order.CreditAmount) {
		go order.Record(false, model.HistoryStepCash, model.PayFeild{
			PayFee:  utils.Violent2String(form.TotalAmount),
			PayType: form.PayType,
			PaidFee: utils.Violent2String(form.FPayAmount),
		})
	} else {
		go order.Record(false, model.HistoryStepCredit, model.PayFeild{
			PayFee:  utils.Violent2String(form.TotalAmount),
			PayType: form.PayType,
			PaidFee: utils.Violent2String(form.FPayAmount),
		})
	}

	// order := logic.NewBatchOrderLogic(c)

	common.Response(c, nil, order)
}

func BatchOrderUpdateStatus(c *gin.Context) {
	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := order.UpdateStatus(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, order)
}

func BatchOrderShared(c *gin.Context) {
	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := order.Shared(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, order)
}

func BatchOrderDetail(c *gin.Context) {

	type PayLoad struct {
		UUID string `json:"uuid"`
	}
	var payLoad PayLoad

	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	if err := order.FromUUID(payLoad.UUID); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, order)
}

func BatchOrderGoodsLatest(c *gin.Context) {
	type Body struct {
		GoodsUUIDList []string `json:"goodsUUIDList"`
	}
	var body Body

	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&body); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := order.FindLatestGoods(body.GoodsUUIDList)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodPriceObjs)
}

func BatchOrderGoodsList(c *gin.Context) {
	type PayLoad struct {
		model.LimitCond
		Status    int32  `json:"status"`
		UserUUID  string `json:"userUUID"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
	}
	var payload PayLoad

	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := order.List(payload.UserUUID, payload.StartTime, payload.EndTime, payload.Status, model.DefaultSetLimitCond(payload.LimitCond))
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodPriceObjs)
}

func BatchGoodsOrderList(c *gin.Context) {
	var payLoad request.BatchGoodsListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}
	logic := logic.NewBatchOrderLogic(c)
	rsp, err := logic.GoodsList(payLoad)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, rsp)
}
