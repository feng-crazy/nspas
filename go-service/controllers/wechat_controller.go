package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
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
	// 获取state参数
	state := ctx.Query("state")
	if state == "" {
		state = "default_state"
	}

	// 生成微信授权URL
	authURL := c.wechatService.GetWeChatOAuthURL(ctx, state)

	ctx.JSON(http.StatusOK, gin.H{
		"url": authURL,
	})
}

// WeChatCallback 处理微信回调
func (c *WeChatController) WeChatCallback(ctx *gin.Context) {
	// 获取code和state参数
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	// 处理微信登录
	user, token, err := c.wechatService.WeChatLogin(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login with WeChat"})
		return
	}

	// 返回登录结果，这里可以根据需要重定向到前端页面
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID.Hex(),
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role,
			"created_at": user.CreatedAt,
		},
		"state": state,
	})
}
