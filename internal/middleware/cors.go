package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Только preflight (OPTIONS) должен устанавливать CORS заголовки
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// ВАЖНО:
		// На обычных запросах НЕ трогаем Access-Control-Allow-Origin,
		// чтобы не было дублирования с backend service через ReverseProxy.
		c.Next()
	}
}
