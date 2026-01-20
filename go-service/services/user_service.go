package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nspas/go-service/config"
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
	// 检查用户是否已存在
	var existingUser models.User
	err := s.col.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, ErrUserExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
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
	_, err = s.col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// 查找用户
	var user models.User
	err := s.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// 生成JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

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
