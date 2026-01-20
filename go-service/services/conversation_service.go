package services

import (
	"context"
	"time"

	"github.com/nspas/go-service/database"
	"github.com/nspas/go-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationService struct {}

func NewConversationService() *ConversationService {
	return &ConversationService{}
}

// CreateConversation 创建对话
func (s *ConversationService) CreateConversation(ctx context.Context, userID primitive.ObjectID, convType models.ConversationType, title string) (*models.Conversation, error) {
	// 创建对话
	conversation := &models.Conversation{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Type:      convType,
		Title:     title,
		Messages:  []models.Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	collection := database.GetCollection("conversations")
	_, err := collection.InsertOne(ctx, conversation)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

// UpdateConversation 更新对话
func (s *ConversationService) UpdateConversation(ctx context.Context, convID primitive.ObjectID, messages []models.Message) (*models.Conversation, error) {
	collection := database.GetCollection("conversations")

	// 更新对话
	update := bson.M{
		"$set": bson.M{
			"messages":   messages,
			"updated_at": time.Now(),
		},
	}

	// 执行更新
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": convID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		return nil, ErrConversationNotFound
	}

	var updatedConv models.Conversation
	err := result.Decode(&updatedConv)
	if err != nil {
		return nil, err
	}

	return &updatedConv, nil
}

// GetConversationByID 根据ID获取对话
func (s *ConversationService) GetConversationByID(ctx context.Context, convID primitive.ObjectID) (*models.Conversation, error) {
	collection := database.GetCollection("conversations")

	var conversation models.Conversation
	err := collection.FindOne(ctx, bson.M{"_id": convID}).Decode(&conversation)
	if err != nil {
		return nil, ErrConversationNotFound
	}

	return &conversation, nil
}

// GetUserConversations 获取用户的所有对话
func (s *ConversationService) GetUserConversations(ctx context.Context, userID primitive.ObjectID) ([]models.Conversation, error) {
	collection := database.GetCollection("conversations")

	// 查询用户的所有对话
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历结果
	var conversations []models.Conversation
	for cursor.Next(ctx) {
		var conv models.Conversation
		if err := cursor.Decode(&conv); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// DeleteConversation 删除对话
func (s *ConversationService) DeleteConversation(ctx context.Context, convID primitive.ObjectID) error {
	collection := database.GetCollection("conversations")

	// 执行删除
	result, err := collection.DeleteOne(ctx, bson.M{"_id": convID})
	if err != nil {
		return err
	}

	// 检查是否删除了记录
	if result.DeletedCount == 0 {
		return ErrConversationNotFound
	}

	return nil
}
