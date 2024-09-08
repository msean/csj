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

func LoggerMiddleware() gin.HandlerFunc {
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
		global.Global.Logger.Info(
			"[response] info",
			zap.Int("status_code", statusCode),
			zap.Float64("duration", time.Now().Sub(_st).Seconds()),
		)
	}
}
