package handler

import (
	"app/global"
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
	var param request.CreateTempBatchOrderParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	logic := logic.NewBatchOrderLogic(c)
	var batchOrder model.BatchOrder
	if batchOrder, err = logic.TempCreate(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	go logic.Record(batchOrder, false, model.HistoryStepOrder, model.PayFeild{})
	common.Response(c, nil, batchOrder)
}

func BatchOrderCreate(c *gin.Context) {
	var param request.CreateBatchOrderParam
	var err error
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchOrderLogic := logic.NewBatchOrderLogic(c)
	var batchOrder model.BatchOrder
	tx := global.Global.DB.Begin()
	if batchOrder, err = batchOrderLogic.Create(tx, param); err != nil {
		tx.Rollback()
		common.Response(c, err, nil)
		return
	}

	orderPayLogic := logic.NewBatchOrderPayLogic(c)
	// orderPayLogic.BatchOrderUUID = batchOrder.UID
	batchOrderPayParam := request.CreateBatchOrderPayParam{
		BatchOrderUUID: batchOrder.UID,
		Amount:         batchOrder.TotalAmount - batchOrder.CreditAmount,
		PayType:        param.PayType,
		CustomerUUID:   batchOrder.UserUUID,
	}
	if _, err = orderPayLogic.Create(tx, batchOrderPayParam, false); err != nil {
		tx.Rollback()
		common.Response(c, err, nil)
		return
	}
	tx.Commit()
	if !common.FloatGreat(0.0, batchOrder.CreditAmount) {
		go batchOrderLogic.Record(batchOrder, false, model.HistoryStepCredit, model.PayFeild{
			PayFee:  param.FPayAmount,
			PayType: param.PayType,
			PaidFee: param.FPayAmount,
		})
	} else {
		go batchOrderLogic.Record(batchOrder, false, model.HistoryStepCash, model.PayFeild{
			PayFee:  batchOrder.TotalAmount - batchOrder.CreditAmount,
			PayType: param.PayType,
			PaidFee: batchOrder.TotalAmount - batchOrder.CreditAmount,
		})
	}

	common.Response(c, nil, batchOrder)
}

func BatchOrderUpdate(c *gin.Context) {
	var param request.UpdateBatchOrderParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	logic := logic.NewBatchOrderLogic(c)
	var batchOrder model.BatchOrder
	if batchOrder, err = logic.Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchOrder)
}

func BatchOrderUpdateStatus(c *gin.Context) {
	var param request.UpdateBatchOrderStatusParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchOrder model.BatchOrder
	if batchOrder, err = logic.NewBatchOrderLogic(c).UpdateStatus(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchOrder)
}

func BatchOrderShared(c *gin.Context) {
	var param request.ShareBatchOrderrParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchOrder model.BatchOrder
	if batchOrder, err = logic.NewBatchOrderLogic(c).Shared(param.UUID); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchOrder)
}

func BatchOrderDetail(c *gin.Context) {
	var param request.BatchOrderDetailParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	logic := logic.NewBatchOrderLogic(c)
	var batchOrder model.BatchOrder
	if batchOrder, err = logic.FromUUID(param.UUID); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchOrder)
}

func BatchOrderGoodsLatest(c *gin.Context) {

	var param request.GoodsLatestOrderParam
	var err error

	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := logic.NewBatchOrderLogic(c).FindLatestGoods(param.GoodsUUIDList)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodPriceObjs)
}

func BatchOrderGoodsList(c *gin.Context) {

	var param request.ListBatchOrderParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodPriceObjs, err := logic.NewBatchOrderLogic(c).List(param)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodPriceObjs)
}
