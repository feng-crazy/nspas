package services

import (
	"context"
	"time"

	"neuro-guide-go-service/database"
	"neuro-guide-go-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService handles user-related business logic
type UserService struct {
	collection *mongo.Collection
}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{
		collection: database.Database.Collection("users"),
	}
}

// CreateUser creates a new user
func (us *UserService) CreateUser(wechatID, nickname, avatar string, isGuest bool) (*models.User, error) {
	user := &models.User{
		ID:        primitive.NewObjectID().Hex(),
		WechatID:  wechatID,
		Nickname:  nickname,
		Avatar:    avatar,
		IsGuest:   isGuest, // 设置是否为游客用户
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return nil, err
	}

	_, err = us.collection.InsertOne(ctx, bson.M{
		"_id":        objID,
		"wechat_id":  user.WechatID,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"is_guest":   user.IsGuest,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (us *UserService) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = us.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByWechatID retrieves a user by their WeChat ID
func (us *UserService) GetUserByWechatID(wechatID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := us.collection.FindOne(ctx, bson.M{"wechat_id": wechatID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates user information
func (us *UserService) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = us.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// GetOrCreateUser gets an existing user or creates a new one
func (us *UserService) GetOrCreateUser(wechatID, nickname, avatar string) (*models.User, error) {
	user, err := us.GetUserByWechatID(wechatID)
	if err == mongo.ErrNoDocuments {
		// User doesn't exist, create new one
		// 默认情况下创建的用户不是游客
		return us.CreateUser(wechatID, nickname, avatar, false)
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// BindPhoneNumber binds a phone number to a user
func (us *UserService) BindPhoneNumber(userID, phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// 更新用户信息，添加手机号（实际项目中应该有更完善的手机号管理）
	update := bson.M{
		"$set": bson.M{
			"phone_number": phoneNumber,
			"is_guest":     false, // 绑定手机号后不再是游客
			"updated_at":   time.Now(),
		},
	}

	_, err = us.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// GetUserByPhoneNumber retrieves a user by their phone number
func (us *UserService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := us.collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateGuestUser creates a new guest user (for development purposes)
func (us *UserService) CreateGuestUser(userID, nickname string) (*models.User, error) {
	// 在实际应用中，游客用户可能有特殊的前缀或格式
	// 这里我们简单地使用传入的ID作为微信ID
	return us.CreateUser(userID, nickname, "", true)
}
