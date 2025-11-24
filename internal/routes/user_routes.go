package routes

import (
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes wires user-service routes to the reverse proxy.
func RegisterUserRoutes(r *gin.Engine, proxy *httputil.ReverseProxy) {
	r.GET("/users/me", proxyHandler(proxy))
	r.GET("/users/:id", proxyHandler(proxy))
	r.POST("/friends/request", proxyHandler(proxy))
	r.GET("/friends/requests/incoming", proxyHandler(proxy))
	r.POST("/friends/requests/:id/accept", proxyHandler(proxy))
	r.POST("/friends/requests/:id/reject", proxyHandler(proxy))
	r.GET("/friends", proxyHandler(proxy))
	r.DELETE("/friends/:id", proxyHandler(proxy))
}
