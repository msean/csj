package appmgr

import (
	"github.com/gin-gonic/gin"
)

type AppUserRouter struct{}

// InitAppUserRouter 初始化 App 注册用户 路由信息
func (s *AppUserRouter) InitAppUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	appUserRouter := Router.Group("appUser")
	{
		appUserRouter.GET("findAppUser", appUserApi.FindAppUser)       // 根据uid获取用户
		appUserRouter.GET("getAppUserList", appUserApi.GetAppUserList) // 获取用户列表
	}
}
