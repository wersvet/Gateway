package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger provides structured logging with request metadata.
func Logger() gin.HandlerFunc {
	logger := log.Default()
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		requestID := GetRequestID(c)

		logger.Printf("status=%d method=%s path=%s latency=%s request_id=%s", status, c.Request.Method, c.Request.URL.Path, latency, requestID)
	}
}
