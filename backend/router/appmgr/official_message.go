package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/middleware"
)

type OfficialMessageRouter struct{}

// InitOfficialMessageRouter 初始化 官方消息 路由信息
func (s *OfficialMessageRouter) InitOfficialMessageRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	officialMessageRouter := Router.Group("officialMessage").Use(middleware.OperationRecord())
	officialMessageRouterWithoutRecord := Router.Group("officialMessage")
	{
		officialMessageRouter.POST("createOfficialMessage", officialMessageApi.CreateOfficialMessage)             // 新建官方消息
		officialMessageRouter.DELETE("deleteOfficialMessage", officialMessageApi.DeleteOfficialMessage)           // 删除官方消息
		officialMessageRouter.DELETE("deleteOfficialMessageByIds", officialMessageApi.DeleteOfficialMessageByIds) // 批量删除官方消息
		officialMessageRouter.PUT("updateOfficialMessage", officialMessageApi.UpdateOfficialMessage)              // 更新官方消息
	}
	{
		officialMessageRouterWithoutRecord.GET("findOfficialMessage", officialMessageApi.FindOfficialMessage)       // 根据uid获取官方消息
		officialMessageRouterWithoutRecord.GET("getOfficialMessageList", officialMessageApi.GetOfficialMessageList) // 获取官方消息列表
	}
}
