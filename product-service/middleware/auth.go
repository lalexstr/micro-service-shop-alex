package middleware

import (
	"net/http"
	"ooolalex/product-service/clients"
	"ooolalex/product-service/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.LoadConfig()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			// Try user_id as fallback
			userIDFloat, ok = claims["user_id"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
				return
			}
		}
		
		userID := uint(userIDFloat)
		c.Set("userID", userID)
		
		// Check role from auth-service
		authClient := clients.NewAuthClient()
		role, err := authClient.GetUserRole(userID)
		if err == nil {
			c.Set("isAdmin", role == "admin")
		} else {
			// Fallback to token claim if auth-service is unavailable
			if isAdmin, ok := claims["is_admin"].(bool); ok {
				c.Set("isAdmin", isAdmin)
			} else {
				c.Set("isAdmin", false)
			}
		}
		
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}
