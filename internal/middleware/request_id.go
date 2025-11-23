package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

// RequestID injects a unique request id into context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set(requestIDKey, id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

// GetRequestID retrieves the request id from context.
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
