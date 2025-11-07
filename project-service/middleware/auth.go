package middleware

import (
	"net/http"
	"ooolalex/project-service/clients"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ContextUserID = "userID"
const ContextIsAdmin = "isAdmin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		// Extract user_id from token (auth-service uses "sub" for user_id)
		var userID uint
		if sub, ok := claims["sub"].(float64); ok {
			userID = uint(sub)
		} else if userIDFloat, ok := claims["user_id"].(float64); ok {
			userID = uint(userIDFloat)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
			return
		}

		c.Set(ContextUserID, userID)

		// Check role from auth-service
		authClient := clients.NewAuthClient()
		role, err := authClient.GetUserRole(userID)
		if err == nil {
			c.Set(ContextIsAdmin, role == "admin")
		} else {
			// Fallback to token claim if auth-service is unavailable
			if isAdmin, ok := claims["is_admin"].(bool); ok {
				c.Set(ContextIsAdmin, isAdmin)
			} else {
				c.Set(ContextIsAdmin, false)
			}
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get(ContextIsAdmin)
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}
