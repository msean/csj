package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func batchOrderPayRouter(g *gin.RouterGroup) {
	batchOrderPayGroup := g.Group("/batch_order_pay", middleware.AuthMiddleware())
	{
		batchOrderPayGroup.POST("/create", BatchOrderPayCreate)
		batchOrderPayGroup.POST("/update", BatchOrderPayUpdate)
		batchOrderPayGroup.POST("/detail", BatchOrderPayDetail)
	}
}

func BatchOrderPayCreate(c *gin.Context) {
	var param request.CreateBatchOrderPayParam
	var err error

	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchOrderPay model.BatchOrderPay
	if batchOrderPay, err = logic.NewBatchOrderPayLogic(c).Create(nil, param, true); err != nil {
		common.Response(c, err, nil)
		return
	}
	// todo记录
	common.Response(c, nil, batchOrderPay)
}

func BatchOrderPayUpdate(c *gin.Context) {
	var param request.UpdateBatchOrderPayParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchOrderPay model.BatchOrderPay
	if batchOrderPay, err = logic.NewBatchOrderPayLogic(c).Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchOrderPay)
}

func BatchOrderPayDetail(c *gin.Context) {
	var param request.GetBatchOrderPayParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchOrderPay model.BatchOrderPay
	if batchOrderPay, err = logic.NewBatchOrderPayLogic(c).FromUUID(param.BatchOrderPayUUID); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchOrderPay)
}
