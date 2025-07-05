package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
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
	var param request.CustomerParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}
	// 增加客户
	customerLogic := logic.NewCustomerLogic(c)
	if param.UID == "" {
		dup, err := customerLogic.ExistName(param.Name)
		if err != nil {
			common.Response(c, err, nil)
			return
		}
		if dup {
			common.Response(c, common.CustomerDuplicateErr, nil)
			return
		}
		var customer model.Customer
		if customer, err = customerLogic.Create(param); err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, customer)
		return
	}

	// 修改客户
	if err = customerLogic.Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	// var rsp response.ListCustomerRsp
	// rsp, err = customerLogic.FromUUID(param.UIDCompatible)
	common.Response(c, err, param)
}

func CustomerList(c *gin.Context) {

	var form request.ListCustomerParam
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewCustomerLogic(c).ListCustomersByOwnerUser(form.SearchKey, form.LimitCond)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, rsp)
}
