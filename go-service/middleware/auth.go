package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger.Info(ctx, "Authentication middleware started")

		// 从Authorization头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn(ctx, "Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logger.Warn(ctx, "Invalid authorization header format", slog.String("header", authHeader))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		logger.Debug(ctx, "Parsing JWT token")

		// 解析token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			logger.Error(ctx, "Failed to parse JWT token", slog.Any("error", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 验证token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户ID存储到上下文
			userID := claims["user_id"].(string)
			email := claims["email"].(string)
			c.Set("user_id", userID)
			c.Set("email", email)
			logger.Info(ctx, "Authentication successful", 
				slog.String("user_id", userID), 
				slog.String("email", email))
			c.Next()
		} else {
			logger.Warn(ctx, "Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
