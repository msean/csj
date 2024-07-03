package csj_customers

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CustomersRouter struct {}

// InitCustomersRouter 初始化 customers表 路由信息
func (s *CustomersRouter) InitCustomersRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	customersRouter := Router.Group("customers").Use(middleware.OperationRecord())
	customersRouterWithoutRecord := Router.Group("customers")
	customersRouterWithoutAuth := PublicRouter.Group("customers")

	var customersApi = v1.ApiGroupApp.Csj_customersApiGroup.CustomersApi
	{
		customersRouter.POST("createCustomers", customersApi.CreateCustomers)   // 新建customers表
		customersRouter.DELETE("deleteCustomers", customersApi.DeleteCustomers) // 删除customers表
		customersRouter.DELETE("deleteCustomersByIds", customersApi.DeleteCustomersByIds) // 批量删除customers表
		customersRouter.PUT("updateCustomers", customersApi.UpdateCustomers)    // 更新customers表
	}
	{
		customersRouterWithoutRecord.GET("findCustomers", customersApi.FindCustomers)        // 根据ID获取customers表
		customersRouterWithoutRecord.GET("getCustomersList", customersApi.GetCustomersList)  // 获取customers表列表
	}
	{
	    customersRouterWithoutAuth.GET("getCustomersPublic", customersApi.GetCustomersPublic)  // 获取customers表列表
	}
}
