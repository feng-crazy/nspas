package controllers

import (
	"fmt"
	"net/http"

	"neuro-guide-go-service/config"
	"neuro-guide-go-service/services"

	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService()
var wechatAuthService *services.WeChatAuthService

// InitUserController initializes the user controller with config
func InitUserController(cfg *config.Config) {
	wechatAuthService = services.NewWeChatAuthService(cfg, userService)
}

// LoginRequest represents a login request
type LoginRequest struct {
	WechatID string `json:"wechat_id" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// WeChatLoginRequest represents a WeChat login request
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// UserLogin handles user login/registration
func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userService.GetOrCreateUser(req.WechatID, req.Nickname, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// WeChatLogin handles WeChat OAuth login
func WeChatLogin(c *gin.Context) {
	var req WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理可能的表单数据（微信网页端登录）
	// 实际项目中应该使用更合适的方式来获取这些数据
	userInfo := make(map[string]interface{})
	if c.Request.Method == "POST" {
		// 简单地从表单中获取数据，实际项目中应该使用更可靠的方式
		if nickname := c.PostForm("nickname"); nickname != "" {
			userInfo["nickname"] = nickname
		}
		if avatar := c.PostForm("avatar"); avatar != "" {
			userInfo["avatar"] = avatar
		}
	}

	user, token, err := wechatAuthService.AuthenticateWithCode(req.Code, &userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("WeChat login failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": gin.H{
				"id":        user.ID,
				"wechat_id": user.WechatID,
				"nickname":  user.Nickname,
				"avatar":    user.Avatar,
				"is_guest":  user.IsGuest,
			},
			"token": token,
		},
	})
}

// BindPhoneNumber handles phone number binding for WeChat users
func BindPhoneNumber(c *gin.Context) {
	// 获取当前用户ID（从JWT token中）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权的用户"})
		return
	}

	// 获取请求参数
	var req struct {
		EncryptedData string `json:"encryptedData" binding:"required"`
		IV            string `json:"iv" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 在实际应用中，这里应该调用微信的解密接口
	// 这里只是示例，直接返回成功

	// 更新用户信息，绑定手机号（实际项目中应该有更完善的手机号管理）
	err := userService.BindPhoneNumber(userID.(string), "13800138000") // 使用虚拟手机号
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("手机号绑定失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "手机号绑定成功",
	})
}

// GetUserProfile handles getting user profile
func GetUserProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile handles updating user profile
type UpdateUserProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func UpdateUserProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
