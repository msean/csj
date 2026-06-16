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

func orderRouter(g *gin.RouterGroup) {
	OrderGroup := g.Group("/batch_order", middleware.AuthMiddleware())
	{
		// 记账单
		OrderGroup.POST("/temp/create", TempOrderCreate)
		// 下单
		OrderGroup.POST("/create", OrderCreate)
		OrderGroup.POST("/update", OrderUpdate)
		OrderGroup.POST("/update_status", OrderUpdateStatus)
		OrderGroup.POST("/detail", OrderDetail)
		OrderGroup.POST("/list", OrderList)
		OrderGroup.POST("/share", OrderShared)
		OrderGroup.POST("/latest_by_goods", OrderGoodsLatest)
		OrderGroup.POST("/credit/list", CreditList)
	}
}

// 码单
func TempOrderCreate(c *gin.Context) {
	order := logic.NewOrderLogic(c)
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

func OrderCreate(c *gin.Context) {
	var form request.OrderReq
	var err error
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	form.OwnerUser = common.GetUserUUID(c)
	order := logic.NewOrderLogic(c)
	err = order.Create(form)
	common.Response(c, err, order)
}

func OrderUpdate(c *gin.Context) {
	var err error
	var form request.OrderReq
	if err = c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewOrderLogic(c)

	err = order.Update(form)
	common.Response(c, err, order)
}

func OrderUpdateStatus(c *gin.Context) {
	order := logic.NewOrderLogic(c)
	var err error
	if err = c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = order.UpdateStatus()
	common.Response(c, err, order)
}

func OrderShared(c *gin.Context) {
	order := logic.NewOrderLogic(c)
	var err error
	if err = c.ShouldBind(&order); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = order.Shared()
	common.Response(c, err, order)
}

func OrderDetail(c *gin.Context) {
	var payLoad request.BatchOrderDetailReq
	var err error
	if err = c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewOrderLogic(c)
	err = order.FromUUID(payLoad.UUID)
	common.Response(c, err, order)
}

func OrderGoodsLatest(c *gin.Context) {

	var body request.BatchOrderGoodsLatest
	var err error
	order := logic.NewOrderLogic(c)
	if err = c.ShouldBind(&body); err != nil {
		common.Response(c, err, nil)
		return
	}

	var goodPriceObjs []response.BatchOrderGoodsOrderRsp
	goodPriceObjs, err = order.FindLatestGoods(body.GoodsUUIDList)
	common.Response(c, err, goodPriceObjs)
}

func OrderList(c *gin.Context) {
	var payload request.BatchOrderListReq
	order := logic.NewOrderLogic(c)
	if err := c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := order.List(payload.UserUUID, payload.StartTime, payload.EndTime, payload.Status, utils.DefaultSetLimitCond(payload.LimitCond))
	common.Response(c, err, goodPriceObjs)
}

func OrderGoodsList(c *gin.Context) {
	var payLoad request.BatchGoodsListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewOrderLogic(c).GoodsList(payLoad)
	common.Response(c, err, rsp)
}

// CreditList 赊欠列表接口
func CreditList(c *gin.Context) {
	var payLoad request.CreditListReq
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewOrderLogic(c).CreditList(payLoad)
	common.Response(c, err, rsp)
}
