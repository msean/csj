package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"

	"github.com/gin-gonic/gin"
)

func memberRouter(g *gin.RouterGroup) {
	group := g.Group("/member", middleware.AuthMiddleware())
	{
		group.POST("/config", MemberConfigHandler)
	}
}

// MemberConfigHandler 返回当前会员收费配置
func MemberConfigHandler(c *gin.Context) {
	configs, err := logic.NewMemberLogic(c).GetMemberConfigs()
	common.Response(c, err, configs)
}
