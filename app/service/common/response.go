package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, err error, data any) {
	if err != nil {
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

// func GetUserUUID(context *gin.Context) string {
// 	// userUUID, _ := c.Value("userUUID").(string)
// 	// return userUUID
// 	// 调试：分别测试两种获取方式
// 	if val, exists := context.Get("userUUID"); exists {
// 		global.Global.Logger.Debug("Get方法获取到:", zap.Any("value", val))
// 	}

//		if val := context.Value("userUUID"); val != nil {
//			global.Global.Logger.Debug("Value方法获取到:", zap.Any("value", val))
//		}
//		return
//	}
func GetUserUUID(c *gin.Context) string {
	// 使用 c.Get() 而不是 c.Value()
	userUUID, exists := c.Get("userUUID")
	if !exists {
		return ""
	}

	// 类型断言，确保安全
	if uuid, ok := userUUID.(string); ok {
		return uuid
	}

	return ""
}
