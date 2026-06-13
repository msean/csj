package handler

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"

	"github.com/gin-gonic/gin"
)

func batchOrderRouter(g *gin.RouterGroup) {
	batchOrderGroup := g.Group("/batch_order", middleware.AuthMiddleware())
	{
		// 记账单
		batchOrderGroup.POST("/temp/create", BatchOrderTempCreate)
		// 下单
		batchOrderGroup.POST("/create", BatchOrderCreate)
		batchOrderGroup.POST("/update", BatchOrderUpdate)
		batchOrderGroup.POST("/update_status", BatchOrderUpdateStatus)
		batchOrderGroup.POST("/detail", BatchOrderDetail)
		batchOrderGroup.POST("/list", BatchOrderList)
		batchOrderGroup.POST("/share", BatchOrderShared)
		batchOrderGroup.POST("/latest_by_goods", BatchOrderGoodsLatest)
		batchOrderGroup.POST("/credit/list", CreditList)
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
	var form request.OrderReq
	var err error
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	err = order.Create(form)
	common.Response(c, err, order)
}

func BatchOrderUpdate(c *gin.Context) {
	var err error
	var form request.OrderReq
	if err = c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)

	err = order.Update(form)
	common.Response(c, err, order)
}

func BatchOrderUpdateStatus(c *gin.Context) {
	order := logic.NewBatchOrderLogic(c)
	var err error
	if err = c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = order.UpdateStatus()
	common.Response(c, err, order)
}

func BatchOrderShared(c *gin.Context) {
	order := logic.NewBatchOrderLogic(c)
	var err error
	if err = c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = order.Shared()
	common.Response(c, err, order)
}

func BatchOrderDetail(c *gin.Context) {
	var payLoad request.BatchOrderDetailReq
	var err error
	if err = c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	err = order.FromUUID(payLoad.UUID)
	common.Response(c, err, order)
}

func BatchOrderGoodsLatest(c *gin.Context) {

	var body request.BatchOrderGoodsLatest
	var err error
	order := logic.NewBatchOrderLogic(c)
	if err = c.ShouldBind(&body); err != nil {
		common.Response(c, err, nil)
		return
	}

	var goodPriceObjs []response.BatchOrderGoodsOrderRsp
	goodPriceObjs, err = order.FindLatestGoods(body.GoodsUUIDList)
	common.Response(c, err, goodPriceObjs)
}

func BatchOrderList(c *gin.Context) {
	var payload request.BatchOrderListReq
	order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := order.List(payload.UserUUID, payload.StartTime, payload.EndTime, payload.Status, utils.DefaultSetLimitCond(payload.LimitCond))
	common.Response(c, err, goodPriceObjs)
}

func BatchGoodsOrderList(c *gin.Context) {
	var payLoad request.BatchGoodsListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewBatchOrderLogic(c).GoodsList(payLoad)
	common.Response(c, err, rsp)
}

// CreditList 赊欠列表接口
func CreditList(c *gin.Context) {
	var payLoad request.CreditListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewBatchOrderLogic(c).CreditList(payLoad)
	common.Response(c, err, rsp)
}
