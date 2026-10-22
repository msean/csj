package handler

import (
	"app/service/handler/middleware"

	"github.com/gin-gonic/gin"
)

func InitEngine(engine *gin.Engine) {
	engine.Use(middleware.LoggerMiddleware())
	engine.StaticFile("/static", "/etc/caishuji/index.html")
	engine.Static("/downloads", "/etc/caishuji")

	apiGroup := engine.Group("/api/csj")
	{
		apiGroup.GET("/healthy", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok!",
			})
		})
		commonRouter(apiGroup)
		userRouter(apiGroup)
		customerRouter(apiGroup)
		goodsRouter(apiGroup)
		batchRouter(apiGroup)
		orderRouter(apiGroup)
		paymentRouter(apiGroup)
		messageRouter(apiGroup)
		memberRouter(apiGroup)
		rechargeRouter(apiGroup)
	}
}
