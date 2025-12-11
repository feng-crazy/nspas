package middleware

import (
	"net/http"
	"strings"

	"neuro-guide-go-service/config"
	"neuro-guide-go-service/services"

	"github.com/gin-gonic/gin"
)

var wechatAuthService *services.WeChatAuthService

// InitAuthMiddleware initializes the auth middleware with config
func InitAuthMiddleware(cfg *config.Config) {
	wechatAuthService = services.NewWeChatAuthService(cfg, services.NewUserService())
}

// AuthMiddleware is a simple authentication middleware
// In production, this should validate JWT tokens or session tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, we'll use a simple header-based auth
		// In production, implement proper JWT or session-based auth
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			// Try to get user_id from query parameter for development
			userID := c.Query("user_id")
			if userID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
				c.Abort()
				return
			}
			c.Set("user_id", userID)
			c.Next()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		user, err := wechatAuthService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Next()
	}
}

// OptionalAuthMiddleware allows requests with or without authentication
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				// Validate token if present
				if user, err := wechatAuthService.ValidateToken(token); err == nil {
					c.Set("user_id", user.ID)
				}
			}
		}

		// Also check query parameter
		if userID := c.Query("user_id"); userID != "" {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}