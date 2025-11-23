package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery ensures panics are returned as JSON errors.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := GetRequestID(c)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":      "internal server error",
			"request_id": requestID,
		})
	})
}
