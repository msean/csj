package biz

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz"
	biz_request "github.com/flipped-aurora/gin-vue-admin/server/model/biz/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomersApi struct{}

var customersService = service.ServiceGroupApp.BizGroup.CustomersService

// CreateCustomers 创建customers表
// @Tags Customers
// @Summary 创建customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csj_customers.Customers true "创建customers表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /customers/createCustomers [post]
func (customersApi *CustomersApi) CreateCustomers(c *gin.Context) {
	var customers biz.Customers
	err := c.ShouldBindJSON(&customers)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := customersService.CreateCustomers(&customers); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteCustomers 删除customers表
// @Tags Customers
// @Summary 删除customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csj_customers.Customers true "删除customers表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /customers/deleteCustomers [delete]
func (customersApi *CustomersApi) DeleteCustomers(c *gin.Context) {
	ID := c.Query("UID")
	if err := customersService.DeleteCustomers(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteCustomersByIds 批量删除customers表
// @Tags Customers
// @Summary 批量删除customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /customers/deleteCustomersByIds [delete]
func (customersApi *CustomersApi) DeleteCustomersByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := customersService.DeleteCustomersByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateCustomers 更新customers表
// @Tags Customers
// @Summary 更新customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body csj_customers.Customers true "更新customers表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /customers/updateCustomers [put]
func (customersApi *CustomersApi) UpdateCustomers(c *gin.Context) {
	var customers biz.Customers
	err := c.ShouldBindJSON(&customers)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := customersService.UpdateCustomers(customers); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindCustomers 用id查询customers表
// @Tags Customers
// @Summary 用id查询customers表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csj_customers.Customers true "用id查询customers表"
// @Success 200 {object} response.Response{data=object{recustomers=csj_customers.Customers},msg=string} "查询成功"
// @Router /customers/findCustomers [get]
func (customersApi *CustomersApi) FindCustomers(c *gin.Context) {
	ID := c.Query("UID")
	if recustomers, err := customersService.GetCustomers(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(recustomers, c)
	}
}

// GetCustomersList 分页获取customers表列表
// @Tags Customers
// @Summary 分页获取customers表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query csj_customersReq.CustomersSearch true "分页获取customers表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /customers/getCustomersList [get]
func (customersApi *CustomersApi) GetCustomersList(c *gin.Context) {
	var pageInfo biz_request.CustomersSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := customersService.GetCustomersInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetCustomersPublic 不需要鉴权的customers表接口
// @Tags Customers
// @Summary 不需要鉴权的customers表接口
// @accept application/json
// @Produce application/json
// @Param data query csj_customersReq.CustomersSearch true "分页获取customers表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /customers/getCustomersPublic [get]
func (customersApi *CustomersApi) GetCustomersPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的customers表接口信息",
	}, "获取成功", c)
}
