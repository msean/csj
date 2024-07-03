package middleware

import (
	"app/global"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		global.GlobalRunTime.Logger.Info(logrus.Fields{
			"method":    c.Request.Method,
			"uri":       c.Request.RequestURI,
			"client_ip": c.ClientIP(),
			"content":   string(b),
			"token":     c.Request.Header.Get("Authorization"),
		})

		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		c.Next()

		statusCode := c.Writer.Status()
		global.GlobalRunTime.Logger.Info(logrus.Fields{
			"status_code": statusCode,
			// "response":    writer.b.String(),
			"duration": time.Now().Sub(_st).Seconds(),
		})
	}
}
