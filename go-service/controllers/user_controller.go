package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
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
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "User registration started")

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn(reqCtx, "Invalid registration request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug(reqCtx, "Processing registration", slog.String("email", req.Email))

	user, err := c.userService.Register(reqCtx, req.Email, req.Password, req.Phone)
	if err != nil {
		if err == services.ErrUserExists {
			logger.Warn(reqCtx, "User already exists", slog.String("email", req.Email))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		logger.Error(reqCtx, "Failed to register user", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	logger.Info(reqCtx, "User registered successfully", slog.String("user_id", user.ID.Hex()))
	ctx.JSON(http.StatusCreated, gin.H{
		"id":         user.ID.Hex(),
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "User login started")

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn(reqCtx, "Invalid login request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug(reqCtx, "Processing login", slog.String("email", req.Email))

	user, token, err := c.userService.Login(reqCtx, req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			logger.Warn(reqCtx, "Invalid login credentials", slog.String("email", req.Email))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		logger.Error(reqCtx, "Failed to login user", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	logger.Info(reqCtx, "User logged in successfully", slog.String("user_id", user.ID.Hex()))
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID.Hex(),
			"email":      user.Email,
			"phone":      user.Phone,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// GetCurrentUser 获取当前用户
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "Get current user started")

	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		logger.Warn(reqCtx, "User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		logger.Warn(reqCtx, "Invalid user ID", slog.String("user_id", userIDStr.(string)))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	logger.Debug(reqCtx, "Getting user by ID", slog.String("user_id", userID.Hex()))

	// 获取用户信息
	user, err := c.userService.GetUserByID(reqCtx, userID)
	if err != nil {
		logger.Error(reqCtx, "Failed to get user", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	logger.Info(reqCtx, "Got current user successfully", slog.String("user_id", user.ID.Hex()))
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.ID.Hex(),
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}
