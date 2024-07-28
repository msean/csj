package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"

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
	orderPay := logic.NewBatchOrderPayLogic(c)
	if err := c.ShouldBind(&orderPay); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := orderPay.Create(nil, true); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, orderPay)
}

func BatchOrderPayUpdate(c *gin.Context) {
	orderPay := logic.NewBatchOrderPayLogic(c)
	if err := c.ShouldBind(&orderPay); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := orderPay.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, orderPay)
}

func BatchOrderPayDetail(c *gin.Context) {
	orderPay := logic.NewBatchOrderPayLogic(c)
	if err := c.ShouldBind(&orderPay); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := orderPay.FromUUID(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, orderPay)
}
