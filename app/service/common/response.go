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
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": 1,
		"msg":  "",
	})
}

func GetUserUUID(c *gin.Context) string {
	userUUID, _ := c.Value("userUUID").(string)
	return userUUID
}
