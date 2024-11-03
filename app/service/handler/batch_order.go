package handler

import (
	"app/global"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"

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
		FPayAmount    float32 `json:"payAmount"`    // 总计
		FCreditAmount float32 `json:"creditAmount"` // 赊欠
		PayType       int32   `json:"payType"`      // 支付方式
	}
	var form Form
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	order.BatchOrder = form.BatchOrder // owner_user重置成了空字符串
	order.OwnerUser = common.GetUserUUID(c)
	order.TotalAmount = order.BatchOrder.SetTotalAmount()
	order.CreditAmount = form.FCreditAmount
	if common.FloatEqual(order.TotalAmount, order.CreditAmount) || common.FloatGreat(order.CreditAmount, order.TotalAmount) {
		order.Status = model.BatchOrderFinish
	} else {
		order.Status = model.BatchOrderedCredit
	}
	tx := global.Global.DB.Begin()
	if err := order.Create(tx); err != nil {
		common.Response(c, err, nil)
		return
	}
	orderPay := logic.NewBatchOrderPayLogic(c)
	orderPay.BatchOrderUUID = order.UID
	orderPay.Amount = order.TotalAmount - order.CreditAmount
	if err := orderPay.Create(tx, false); err != nil {
		common.Response(c, err, nil)
		return
	}
	tx.Commit()
	if common.FloatGreat(0.0, form.FCreditAmount) {
		go order.Record(false, model.HistoryStepCash, model.PayFeild{
			PayFee:  order.TotalAmount - order.CreditAmount,
			PayType: form.PayType,
			PaidFee: order.TotalAmount - order.CreditAmount,
		})
	} else {
		go order.Record(false, model.HistoryStepCredit, model.PayFeild{
			PayFee:  order.TotalAmount - order.CreditAmount,
			PayType: form.PayType,
			PaidFee: order.TotalAmount - order.CreditAmount,
		})
	}

	common.Response(c, nil, order)
}

func BatchOrderUpdate(c *gin.Context) {
	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := order.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
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
	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := order.FromUUID(); err != nil {
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
