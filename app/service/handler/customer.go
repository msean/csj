package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"

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
	if err := c.ShouldBind(&customerLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	// 增加客户
	if customerLogic.UID == "" {
		dup, err := customerLogic.Check()
		if err != nil {
			common.Response(c, err, nil)
			return
		}
		if dup {
			common.Response(c, common.CustomerDuplicateErr, nil)
			return
		}
		if err := customerLogic.Create(); err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, customerLogic)
		return
	}

	// 修改客户
	if err := customerLogic.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, customerLogic)
}

func CustomerList(c *gin.Context) {

	type Form struct {
		model.LimitCond
		SearchKey string `json:"searchName"`
	}
	var form Form
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	var _customers []logic.CustomerLogic
	var err error
	_customers, err = logic.NewCustomerLogic(c).ListCustomersByOwnerUser(form.SearchKey, form.LimitCond)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, _customers)
}
