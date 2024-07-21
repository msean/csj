package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"time"

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
		PayAmount float32 `json:"pay_amount"`
		PayType   float32 `json:"pay_type"`
	}
	var form Form
	// order := logic.NewBatchOrderLogic(c)
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	order := logic.NewBatchOrderLogic(c)
	order.BatchOrder = form.BatchOrder
	if err := order.Create(); err != nil {
		common.Response(c, err, nil)
		return
	}
	go order.Record(false, model.HistoryStepOrder, model.PayFeild{})
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
		Status    int32     `json:"status"`
		UserUUID  string    `json:"userUUID"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
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
