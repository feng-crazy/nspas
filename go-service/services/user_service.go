package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
	"github.com/nspas/go-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	cfg     *config.Config
	db      *mongo.Database
	col     *mongo.Collection
}

func NewUserService(cfg *config.Config, db *mongo.Database) *UserService {
	return &UserService{
		cfg:     cfg,
		db:      db,
		col:     db.Collection("users"),
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, email, password, phone string) (*models.User, error) {
	logger.Info(ctx, "User registration started", slog.String("email", email))

	// 检查用户是否已存在
	var existingUser models.User
	err := s.col.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		logger.Warn(ctx, "User already exists", slog.String("email", email))
		return nil, ErrUserExists
	}

	// 加密密码
	logger.Debug(ctx, "Hashing password for user", slog.String("email", email))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "Failed to hash password", slog.Any("error", err))
		return nil, err
	}

	// 创建新用户
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Phone:     phone,
		Password:  string(hashedPassword),
		Role:      models.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	logger.Debug(ctx, "Saving new user to database", slog.String("email", email))
	_, err = s.col.InsertOne(ctx, user)
	if err != nil {
		logger.Error(ctx, "Failed to save user to database", slog.Any("error", err))
		return nil, err
	}

	logger.Info(ctx, "User registered successfully", slog.String("user_id", user.ID.Hex()))
	return user, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	logger.Info(ctx, "User login started", slog.String("email", email))

	// 查找用户
	var user models.User
	err := s.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		logger.Warn(ctx, "User not found", slog.String("email", email))
		return nil, "", ErrInvalidCredentials
	}

	// 验证密码
	logger.Debug(ctx, "Verifying password for user", slog.String("email", email))
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Warn(ctx, "Invalid password", slog.String("email", email))
		return nil, "", ErrInvalidCredentials
	}

	// 生成JWT token
	logger.Debug(ctx, "Generating JWT token for user", slog.String("email", email))
	token, err := s.generateToken(user)
	if err != nil {
		logger.Error(ctx, "Failed to generate JWT token", slog.Any("error", err))
		return nil, "", err
	}

	logger.Info(ctx, "User logged in successfully", slog.String("user_id", user.ID.Hex()))
	return &user, token, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	logger.Info(ctx, "Get user by ID started", slog.String("user_id", id.Hex()))

	var user models.User
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Warn(ctx, "User not found by ID", slog.String("user_id", id.Hex()))
		}
		logger.Error(ctx, "Failed to get user by ID", slog.Any("error", err))
		return nil, err
	}

	logger.Info(ctx, "Got user by ID successfully", slog.String("user_id", user.ID.Hex()))
	return &user, nil
}

// generateToken 生成JWT token
func (s *UserService) generateToken(user models.User) (string, error) {
	// 设置token过期时间
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWT.Expires) * time.Hour)

	// 创建claims
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
