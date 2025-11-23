package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	userIDKey   = "user_id"
	usernameKey = "username"
)

// JWTAuth validates tokens and injects user claims into the context.
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if isPublicEndpoint(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		tokenString := extractTokenFromHeaderOrQuery(c)
		if tokenString == "" {
			unauthorized(c)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			unauthorized(c)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthorized(c)
			return
		}

		userID1, _ := claims["user_id"].(float64)
		userID := strconv.Itoa(int(userID1))
		username, _ := claims["username"].(string)
		if userID == "" || username == "" {
			unauthorized(c)
			return
		}

		c.Set(userIDKey, userID)
		c.Set(usernameKey, username)
		c.Request.Header.Set("X-User-ID", userID)
		c.Request.Header.Set("X-Username", username)
		c.Request.Header.Set("Authorization", "Bearer "+tokenString)

		c.Next()
	}
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

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func isPublicEndpoint(method, path string) bool {
	// allow any /auth/login*
	if method == http.MethodPost && strings.HasPrefix(path, "/auth/login") {
		return true
	}

	// allow any /auth/register*
	if method == http.MethodPost && strings.HasPrefix(path, "/auth/register") {
		return true
	}

	return false
}

func extractTokenFromHeaderOrQuery(c *gin.Context) string {
	// 1) сначала пробуем взять токен из заголовка Authorization
	header := c.GetHeader("Authorization")
	if header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return strings.TrimSpace(parts[1])
		}
	}

	// 2) если нет — пробуем из query параметра: ws://.../ws/chats/1?token=JWT
	token := c.Query("token")
	if token != "" {
		return token
	}

	return ""
}
