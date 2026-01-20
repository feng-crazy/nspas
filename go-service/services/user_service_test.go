package services

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository 模拟用户存储库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	// 创建mock依赖
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:  "test-secret",
			Expires: 24,
		},
	}

	// 创建用户服务
	service := NewUserService(cfg)

	// 测试用例1：成功注册
	t.Run("success", func(t *testing.T) {
		// 设置mock期望
		mockRepo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(nil, ErrNotFound)
		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

		// 执行测试
		user, err := service.Register(context.Background(), "test@example.com", "password123", "1234567890")

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, models.RoleUser, user.Role)
		mockRepo.AssertExpectations(t)
	})

	// 测试用例2：用户已存在
	t.Run("user_exists", func(t *testing.T) {
		// 设置mock期望
		existingUser := &models.User{
			Email: "test@example.com",
		}
		mockRepo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

		// 执行测试
		user, err := service.Register(context.Background(), "test@example.com", "password123", "1234567890")

		// 验证结果
		assert.Error(t, err)
		assert.Equal(t, ErrUserExists, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	// 创建mock依赖
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:  "test-secret",
			Expires: 24,
		},
	}

	// 创建用户服务
	service := NewUserService(cfg)

	// 测试用例1：成功登录
	t.Run("success", func(t *testing.T) {
		// 设置mock期望
		user := &models.User{
			ID:       primitive.NewObjectID(),
			Email:    "test@example.com",
			Password: "$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW", // password123
			Role:     models.RoleUser,
		}
		mockRepo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

		// 执行测试
		resultUser, token, err := service.Login(context.Background(), "test@example.com", "password123")

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, resultUser)
		assert.Equal(t, user.ID, resultUser.ID)
		assert.NotEmpty(t, token)

		// 验证token是否有效
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})

	// 测试用例2：密码错误
	t.Run("invalid_password", func(t *testing.T) {
		// 设置mock期望
		user := &models.User{
			Email:    "test@example.com",
			Password: "$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW", // password123
		}
		mockRepo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(user, nil)

		// 执行测试
		resultUser, token, err := service.Login(context.Background(), "test@example.com", "wrongpassword")

		// 验证结果
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCredentials, err)
		assert.Nil(t, resultUser)
		assert.Empty(t, token)
	})
}
