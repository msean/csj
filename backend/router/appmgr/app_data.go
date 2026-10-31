package appmgr

import "github.com/gin-gonic/gin"

type AppDataRouter struct{}

func (r *AppDataRouter) InitAppDataRouter(private *gin.RouterGroup) {
	group := private.Group("appData")
	group.GET(":resource/list", appDataApi.List)
	group.GET(":resource/detail", appDataApi.Detail)
}
