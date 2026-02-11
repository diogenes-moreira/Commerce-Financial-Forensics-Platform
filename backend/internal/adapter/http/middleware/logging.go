package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logging(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		entry := log.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"latency_ms": time.Since(start).Milliseconds(),
			"ip":         c.ClientIP(),
		})

		if reqID, exists := c.Get("request_id"); exists {
			entry = entry.WithField("request_id", reqID)
		}

		if c.Writer.Status() >= 500 {
			entry.Error("Server error")
		} else {
			entry.Info("Request handled")
		}
	}
}
