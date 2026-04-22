package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BotResponse struct {
	Code int         `json:"code"` // 自定义状态码
	Msg  string      `json:"msg"`  // 消息说明
	Data interface{} `json:"data"` // 返回数据
}

func BotResult(c *gin.Context, httpCode int, code int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func BotSuccess(c *gin.Context, data interface{}) {
	BotResult(c, http.StatusOK, 0, "success", data)
}

// 成功 + 自定义消息
func BotSuccessMsg(c *gin.Context, msg string, data interface{}) {
	BotResult(c, http.StatusOK, 0, msg, data)
}

// 参数错误
func BotBadRequest(c *gin.Context, msg string) {
	BotResult(c, http.StatusBadRequest, 400, msg, nil)
}

// 未授权
func BotUnauthorized(c *gin.Context, msg string) {
	BotResult(c, http.StatusUnauthorized, 401, msg, nil)
}

// 资源不存在
func BotNotFound(c *gin.Context, msg string) {
	BotResult(c, http.StatusNotFound, 404, msg, nil)
}

// 服务器内部错误
func BotFail(c *gin.Context, msg string) {
	BotResult(c, http.StatusInternalServerError, 500, msg, nil)
}
