package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"

	"github.com/gin-gonic/gin"
)

func orderPayRouter(g *gin.RouterGroup) {
	orderPayGroup := g.Group("/batch_order_pay", middleware.AuthMiddleware())
	{
		orderPayGroup.POST("/create", OrderPayCreate)
		orderPayGroup.POST("/update", OrderPayUpdate)
		orderPayGroup.POST("/detail", OrderPayDetail)
	}
}

func OrderPayCreate(c *gin.Context) {
	orderPay := logic.NewOrderPayLogic(c)
	if err := c.ShouldBind(&orderPay); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := orderPay.Create(nil, true); err != nil {
		common.Response(c, err, nil)
		return
	}
	// todo记录
	common.Response(c, nil, orderPay)
}

func OrderPayUpdate(c *gin.Context) {
	orderPay := logic.NewOrderPayLogic(c)
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

func OrderPayDetail(c *gin.Context) {
	orderPay := logic.NewOrderPayLogic(c)
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
