package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/middleware"
)

type MemberConfigRouter struct{}

// InitMemberConfigRouter 初始化 会员收费配置 路由信息
func (s *MemberConfigRouter) InitMemberConfigRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	memberConfigRouter := Router.Group("memberConfig").Use(middleware.OperationRecord())
	memberConfigRouterWithoutRecord := Router.Group("memberConfig")
	{
		memberConfigRouter.POST("createMemberConfig", memberConfigApi.CreateMemberConfig)             // 新建会员收费配置
		memberConfigRouter.DELETE("deleteMemberConfig", memberConfigApi.DeleteMemberConfig)           // 删除会员收费配置
		memberConfigRouter.DELETE("deleteMemberConfigByIds", memberConfigApi.DeleteMemberConfigByIds) // 批量删除会员收费配置
		memberConfigRouter.PUT("updateMemberConfig", memberConfigApi.UpdateMemberConfig)              // 更新会员收费配置
	}
	{
		memberConfigRouterWithoutRecord.GET("findMemberConfig", memberConfigApi.FindMemberConfig)       // 根据uid获取会员收费配置
		memberConfigRouterWithoutRecord.GET("getMemberConfigList", memberConfigApi.GetMemberConfigList) // 获取会员收费配置列表
	}
}
