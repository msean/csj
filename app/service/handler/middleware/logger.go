package middleware

import (
	"app/global"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.b.Write(b)
	return w.ResponseWriter.Write(b)
}

func shouldLogResponse() bool {
	if global.Global.Env() == "test" || global.Global.Env() == "debug" {
		return true
	}
	return false
}

// respMaxLen 响应体日志最大长度，超出部分截断，避免日志爆炸
const respMaxLen = 8192

func LoggerMiddleware() gin.HandlerFunc {
	logResp := shouldLogResponse()
	return func(c *gin.Context) {
		_st := time.Now()
		b, _ := c.GetRawData()
		global.Global.Logger.Info(
			"[request] info",
			zap.String("method", c.Request.Method),
			zap.String("uri", c.Request.RequestURI),
			zap.String("client_ip", c.ClientIP()),
			zap.String("content", string(b)),
			zap.String("token", c.Request.Header.Get("Authorization")),
		)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		c.Next()

		statusCode := c.Writer.Status()
		if !logResp {
			global.Global.Logger.Info(
				"[response] info",
				zap.Int("status_code", statusCode),
				zap.Float64("duration", time.Since(_st).Seconds()),
			)
			return
		}

		respBody := writer.b.String()
		if len(respBody) > respMaxLen {
			respBody = respBody[:respMaxLen] + "...(truncated)"
		}
		global.Global.Logger.Info(
			"[response] info",
			zap.Int("status_code", statusCode),
			zap.Float64("duration", time.Since(_st).Seconds()),
			zap.String("body", respBody),
		)
	}
}
