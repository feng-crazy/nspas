package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
	"github.com/nspas/go-service/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type WeChatController struct {
	wechatService *services.WeChatService
	userService   *services.UserService
}

func NewWeChatController(cfg *config.Config, db *mongo.Database) *WeChatController {
	return &WeChatController{
		wechatService: services.NewWeChatService(cfg),
		userService:   services.NewUserService(cfg, db),
	}
}

// GetWeChatAuthURL 获取微信授权URL
func (c *WeChatController) GetWeChatAuthURL(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "Get WeChat auth URL request started")

	// 获取state参数
	state := ctx.Query("state")
	if state == "" {
		state = "default_state"
		logger.Debug(reqCtx, "No state provided, using default", slog.String("default_state", state))
	}

	// 生成微信授权URL
	authURL := c.wechatService.GetWeChatOAuthURL(reqCtx, state)

	logger.Info(reqCtx, "WeChat auth URL generated successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"url": authURL,
	})
}

// WeChatCallback 处理微信回调
func (c *WeChatController) WeChatCallback(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "WeChat callback request started")

	// 获取code和state参数
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		logger.Warn(reqCtx, "WeChat callback missing code parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	logger.Debug(reqCtx, "Processing WeChat callback",
		slog.String("code", code),
		slog.String("state", state))

	// 处理微信登录
	user, token, err := c.wechatService.WeChatLogin(reqCtx, code)
	if err != nil {
		logger.Error(reqCtx, "Failed to login with WeChat", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login with WeChat"})
		return
	}

	logger.Info(reqCtx, "WeChat login successful", slog.String("user_id", user.ID.Hex()))
	// 返回登录结果，这里可以根据需要重定向到前端页面
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID.Hex(),
			"email":      user.Email,
			"phone":      user.Phone,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
		"state": state,
	})
}
