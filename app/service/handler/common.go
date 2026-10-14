package handler

import (
	"app/service/common"
	"app/service/logic"

	"github.com/gin-gonic/gin"
)

func commonRouter(g *gin.RouterGroup) {
	group := g.Group("/common")
	{
		group.POST("/init", InitHandler)
	}
}

func InitHandler(c *gin.Context) {
	data, err := logic.NewCommonLogic(c).Init()
	common.Response(c, err, data)
}
