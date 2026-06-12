package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func customerRouter(g *gin.RouterGroup) {
	group := g.Group("/customer", middleware.AuthMiddleware())
	{
		group.POST("/save", CustomerSave)
		group.POST("/list", CustomerList)
	}
}

func CustomerSave(c *gin.Context) {
	customerLogic := logic.NewCustomerLogic(c)
	var err error
	if err = c.ShouldBind(&customerLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	// 增加客户
	if customerLogic.UID == "" {
		err = customerLogic.Create()
		common.Response(c, err, customerLogic)
		return
	}

	// 修改客户
	err = customerLogic.Update()
	common.Response(c, err, customerLogic)
}

func CustomerList(c *gin.Context) {
	var form request.CustomerListReq
	var err error
	if err = c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}
	var _customers []logic.CustomerLogic

	_customers, err = logic.NewCustomerLogic(c).ListCustomersByOwnerUser(form)
	common.Response(c, err, _customers)
}
