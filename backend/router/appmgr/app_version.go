package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/middleware"
)

type AppVersionRouter struct{}

// InitAppVersionRouter 初始化 App版本 路由信息
func (s *AppVersionRouter) InitAppVersionRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	appVersionRouter := Router.Group("appVersion").Use(middleware.OperationRecord())
	appVersionRouterWithoutRecord := Router.Group("appVersion")
	{
		appVersionRouter.POST("createAppVersion", appVersionApi.CreateAppVersion)             // 新建App版本
		appVersionRouter.DELETE("deleteAppVersion", appVersionApi.DeleteAppVersion)           // 删除App版本
		appVersionRouter.DELETE("deleteAppVersionByIds", appVersionApi.DeleteAppVersionByIds) // 批量删除App版本
		appVersionRouter.PUT("updateAppVersion", appVersionApi.UpdateAppVersion)              // 更新App版本
	}
	{
		appVersionRouterWithoutRecord.GET("findAppVersion", appVersionApi.FindAppVersion)       // 根据ID获取App版本
		appVersionRouterWithoutRecord.GET("getAppVersionList", appVersionApi.GetAppVersionList) // 获取App版本列表
	}
}
