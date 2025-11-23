package routes

import (
	"net/http/httputil"

	"TEST/internal/proxy"
	"github.com/gin-gonic/gin"
)

// RegisterChatRoutes wires chat-service routes and websocket proxy.
func RegisterChatRoutes(r *gin.Engine, proxyHandler *httputil.ReverseProxy, wsProxy *proxy.WebsocketProxy) {
	r.GET("/chats", proxyHandlerFunc(proxyHandler))
	r.POST("/chats/start", proxyHandlerFunc(proxyHandler))
	r.GET("/chats/:chat_id/messages", proxyHandlerFunc(proxyHandler))
	r.POST("/chats/:chat_id/messages", proxyHandlerFunc(proxyHandler))
	r.DELETE("/chats/:chat_id/messages/:message_id/me", proxyHandlerFunc(proxyHandler))
	r.DELETE("/chats/:chat_id/messages/:message_id/all", proxyHandlerFunc(proxyHandler))
	r.DELETE("/chats/:chat_id/me", proxyHandlerFunc(proxyHandler))

	r.GET("/ws/chats/:chat_id", wsProxy.Handle)
}

func proxyHandlerFunc(p *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		p.ServeHTTP(c.Writer, c.Request)
	}
}
