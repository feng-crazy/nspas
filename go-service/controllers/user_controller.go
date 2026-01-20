package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(cfg *config.Config, db *mongo.Database) *UserController {
	return &UserController{
		userService: services.NewUserService(cfg, db),
	}
}

// RegisterRequest 注册请求
 type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

// LoginRequest 登录请求
 type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.Register(ctx, req.Email, req.Password, req.Phone)
	if err != nil {
		if err == services.ErrUserExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":        user.ID.Hex(),
		"email":     user.Email,
		"phone":     user.Phone,
		"role":      user.Role,
		"created_at": user.CreatedAt,
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID.Hex(),
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// GetCurrentUser 获取当前用户
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":        user.ID.Hex(),
		"email":     user.Email,
		"phone":     user.Phone,
		"role":      user.Role,
		"created_at": user.CreatedAt,
	})
}
