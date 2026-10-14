package common

import (
	"fmt"
	"net/http"

	"app/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Response(c *gin.Context, err error, data any) {
	if err != nil {
		global.Global.Logger.Error("API response error",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": 1,
		"msg":  "",
	})
}

func IllegalResponse(c *gin.Context, err error, data any) {
	if err != nil {
		global.Global.Logger.Warn("API illegal response",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
			"code": 4001,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": 1,
		"msg":  "",
	})
}

func GetUserUUID(c *gin.Context) string {
	// 使用 c.Get() 而不是 c.Value()
	userUUID, exists := c.Get("userUUID")
	fmt.Println(">>>>>>>>GetUserUUID", userUUID)
	if !exists {
		return ""
	}

	// 类型断言，确保安全
	if uuid, ok := userUUID.(string); ok {
		return uuid
	}

	return ""
}
