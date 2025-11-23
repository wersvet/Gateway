package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebsocketProxy proxies websocket traffic to the chat service while appending the JWT token.
type WebsocketProxy struct {
	target *url.URL
}

// NewWebsocketProxy creates a proxy for websocket endpoints targeting the chat service base URL.
func NewWebsocketProxy(target string) (*WebsocketProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &WebsocketProxy{target: targetURL}, nil
}

// Handle proxies the incoming websocket request to the backend chat service.
func (p *WebsocketProxy) Handle(c *gin.Context) {
	token := extractToken(c.GetHeader("Authorization"))
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	backendURL := *p.target
	backendURL.Scheme = convertScheme(backendURL.Scheme)
	backendURL.Path = fmt.Sprintf("/ws/chats/%s", c.Param("chat_id"))

	q := backendURL.Query()
	q.Set("token", token)
	backendURL.RawQuery = q.Encode()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to upgrade websocket"})
		return
	}
	defer clientConn.Close()

	dialer := websocket.Dialer{}
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	if userID := c.GetString("user_id"); userID != "" {
		headers.Set("X-User-ID", userID)
	}
	if username := c.GetString("username"); username != "" {
		headers.Set("X-Username", username)
	}

	backendConn, _, err := dialer.Dial(backendURL.String(), headers)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "failed to connect to chat service"})
		return
	}
	defer backendConn.Close()

	errCh := make(chan error, 2)

	go proxyMessages(clientConn, backendConn, errCh)
	go proxyMessages(backendConn, clientConn, errCh)

	if err := <-errCh; err != nil {
		c.Error(err) // logged by gin
	}
}

func proxyMessages(src, dst *websocket.Conn, errCh chan<- error) {
	for {
		msgType, msg, err := src.ReadMessage()
		if err != nil {
			errCh <- err
			return
		}
		if err := dst.WriteMessage(msgType, msg); err != nil {
			errCh <- err
			return
		}
	}
}

func convertScheme(scheme string) string {
	if strings.EqualFold(scheme, "https") {
		return "wss"
	}
	return "ws"
}

func extractToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
