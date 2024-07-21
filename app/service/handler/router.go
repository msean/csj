package handler

import (
	"app/service/handler/middleware"

	"github.com/gin-gonic/gin"
)

func InitEngine(engine *gin.Engine) {
	engine.Use(middleware.LoggerMiddleware())
	engine.StaticFile("/static", "/etc/caishuji/index.html")
	apiGroup := engine.Group("/api/csj")
	{
		apiGroup.GET("/healthy", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok!",
			})
		})
		userRouter(apiGroup)
		customerRouter(apiGroup)
		goodsRouter(apiGroup)
		batchRouter(apiGroup)
		batchOrderRouter(apiGroup)
		batchOrderPayRouter(apiGroup)
	}
}
