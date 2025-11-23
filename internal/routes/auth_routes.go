package routes

import (
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes wires auth-service routes to the reverse proxy.
func RegisterAuthRoutes(r *gin.Engine, proxy *httputil.ReverseProxy) {
	r.POST("/auth/register", proxyHandler(proxy))
	r.POST("/auth/login", proxyHandler(proxy))
	r.GET("/auth/validate", proxyHandler(proxy))
	r.GET("/auth/user/:id", proxyHandler(proxy))
}

func proxyHandler(p *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		p.ServeHTTP(c.Writer, c.Request)
	}
}
