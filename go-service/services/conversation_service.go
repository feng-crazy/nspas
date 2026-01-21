package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/nspas/go-service/database"
	"github.com/nspas/go-service/logger"
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
	logger.Info(ctx, "Creating conversation started",
		slog.String("user_id", userID.Hex()),
		slog.String("type", string(convType)),
		slog.String("title", title))

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
	logger.Debug(ctx, "Saving conversation to database",
		slog.String("conversation_id", conversation.ID.Hex()))
	_, err := collection.InsertOne(ctx, conversation)
	if err != nil {
		logger.Error(ctx, "Failed to save conversation to database",
			slog.String("conversation_id", conversation.ID.Hex()),
			slog.Any("error", err))
		return nil, err
	}

	logger.Info(ctx, "Conversation created successfully",
		slog.String("conversation_id", conversation.ID.Hex()))
	return conversation, nil
}

// UpdateConversation 更新对话
func (s *ConversationService) UpdateConversation(ctx context.Context, convID primitive.ObjectID, messages []models.Message) (*models.Conversation, error) {
	logger.Info(ctx, "Updating conversation started",
		slog.String("conversation_id", convID.Hex()),
		slog.Int("message_count", len(messages)))

	collection := database.GetCollection("conversations")

	// 更新对话
	update := bson.M{
		"$set": bson.M{
			"messages":   messages,
			"updated_at": time.Now(),
		},
	}

	// 执行更新
	logger.Debug(ctx, "Executing conversation update",
		slog.String("conversation_id", convID.Hex()))
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": convID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		logger.Error(ctx, "Conversation not found for update",
			slog.String("conversation_id", convID.Hex()))
		return nil, ErrConversationNotFound
	}

	var updatedConv models.Conversation
	err := result.Decode(&updatedConv)
	if err != nil {
		logger.Error(ctx, "Failed to decode updated conversation",
			slog.String("conversation_id", convID.Hex()),
			slog.Any("error", err))
		return nil, err
	}

	logger.Info(ctx, "Conversation updated successfully",
		slog.String("conversation_id", updatedConv.ID.Hex()))
	return &updatedConv, nil
}

// GetConversationByID 根据ID获取对话
func (s *ConversationService) GetConversationByID(ctx context.Context, convID primitive.ObjectID) (*models.Conversation, error) {
	logger.Info(ctx, "Getting conversation by ID started",
		slog.String("conversation_id", convID.Hex()))

	collection := database.GetCollection("conversations")

	var conversation models.Conversation
	err := collection.FindOne(ctx, bson.M{"_id": convID}).Decode(&conversation)
	if err != nil {
		logger.Error(ctx, "Conversation not found",
			slog.String("conversation_id", convID.Hex()))
		return nil, ErrConversationNotFound
	}

	logger.Info(ctx, "Got conversation by ID successfully",
		slog.String("conversation_id", conversation.ID.Hex()))
	return &conversation, nil
}

// GetUserConversations 获取用户的所有对话
func (s *ConversationService) GetUserConversations(ctx context.Context, userID primitive.ObjectID) ([]models.Conversation, error) {
	logger.Info(ctx, "Getting user conversations started",
		slog.String("user_id", userID.Hex()))

	collection := database.GetCollection("conversations")

	// 查询用户的所有对话
	logger.Debug(ctx, "Executing user conversations query",
		slog.String("user_id", userID.Hex()))
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Error(ctx, "Failed to query user conversations",
			slog.String("user_id", userID.Hex()),
			slog.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历结果
	var conversations []models.Conversation
	for cursor.Next(ctx) {
		var conv models.Conversation
		if err := cursor.Decode(&conv); err != nil {
			logger.Error(ctx, "Failed to decode conversation",
				slog.String("user_id", userID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	logger.Info(ctx, "Got user conversations successfully",
		slog.String("user_id", userID.Hex()),
		slog.Int("conversation_count", len(conversations)))
	return conversations, nil
}

// DeleteConversation 删除对话
func (s *ConversationService) DeleteConversation(ctx context.Context, convID primitive.ObjectID) error {
	logger.Info(ctx, "Deleting conversation started",
		slog.String("conversation_id", convID.Hex()))

	collection := database.GetCollection("conversations")

	// 执行删除
	logger.Debug(ctx, "Executing conversation deletion",
		slog.String("conversation_id", convID.Hex()))
	result, err := collection.DeleteOne(ctx, bson.M{"_id": convID})
	if err != nil {
		logger.Error(ctx, "Failed to delete conversation",
			slog.String("conversation_id", convID.Hex()),
			slog.Any("error", err))
		return err
	}

	// 检查是否删除了记录
	if result.DeletedCount == 0 {
		logger.Error(ctx, "Conversation not found for deletion",
			slog.String("conversation_id", convID.Hex()))
		return ErrConversationNotFound
	}

	logger.Info(ctx, "Conversation deleted successfully",
		slog.String("conversation_id", convID.Hex()))
	return nil
}
