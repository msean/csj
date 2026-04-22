package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/middleware"
	"github.com/msean/csj/backend/plugin/announcement/router"
)

func Router(engine *gin.Engine) {
	public := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	private := engine.Group(global.GVA_CONFIG.System.RouterPrefix).Group("")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	router.Router.Info.Init(public, private)
}
