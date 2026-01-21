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
)

type ToolService struct{}

func NewToolService() *ToolService {
	return &ToolService{}
}

// SaveTool 保存工具
func (s *ToolService) SaveTool(ctx context.Context, userID primitive.ObjectID, name, description, htmlContent string, convID primitive.ObjectID) (*models.Tool, error) {
	logger.Info(ctx, "Saving tool started",
		slog.String("user_id", userID.Hex()),
		slog.String("name", name),
		slog.String("conversation_id", convID.Hex()))

	// 创建工具
	tool := &models.Tool{
		ID:             primitive.NewObjectID(),
		UserID:         userID,
		Name:           name,
		Description:    description,
		HTMLContent:    htmlContent,
		ConversationID: convID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 保存到数据库
	collection := database.GetCollection("tools")
	logger.Debug(ctx, "Saving tool to database",
		slog.String("tool_id", tool.ID.Hex()))
	_, err := collection.InsertOne(ctx, tool)
	if err != nil {
		logger.Error(ctx, "Failed to save tool to database",
			slog.String("tool_id", tool.ID.Hex()),
			slog.Any("error", err))
		return nil, err
	}

	logger.Info(ctx, "Tool saved successfully",
		slog.String("tool_id", tool.ID.Hex()))
	return tool, nil
}

// GetToolByID 根据ID获取工具
func (s *ToolService) GetToolByID(ctx context.Context, toolID primitive.ObjectID) (*models.Tool, error) {
	logger.Info(ctx, "Getting tool by ID started",
		slog.String("tool_id", toolID.Hex()))

	collection := database.GetCollection("tools")

	var tool models.Tool
	err := collection.FindOne(ctx, bson.M{"_id": toolID}).Decode(&tool)
	if err != nil {
		logger.Error(ctx, "Tool not found",
			slog.String("tool_id", toolID.Hex()))
		return nil, ErrToolNotFound
	}

	logger.Info(ctx, "Got tool by ID successfully",
		slog.String("tool_id", tool.ID.Hex()))
	return &tool, nil
}

// GetUserTools 获取用户的所有工具
func (s *ToolService) GetUserTools(ctx context.Context, userID primitive.ObjectID) ([]models.Tool, error) {
	logger.Info(ctx, "Getting user tools started",
		slog.String("user_id", userID.Hex()))

	collection := database.GetCollection("tools")

	// 查询用户的所有工具
	logger.Debug(ctx, "Executing user tools query",
		slog.String("user_id", userID.Hex()))
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Error(ctx, "Failed to query user tools",
			slog.String("user_id", userID.Hex()),
			slog.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历结果
	var tools []models.Tool
	for cursor.Next(ctx) {
		var tool models.Tool
		if err := cursor.Decode(&tool); err != nil {
			logger.Error(ctx, "Failed to decode tool",
				slog.String("user_id", userID.Hex()),
				slog.Any("error", err))
			return nil, err
		}
		tools = append(tools, tool)
	}

	logger.Info(ctx, "Got user tools successfully",
		slog.String("user_id", userID.Hex()),
		slog.Int("tool_count", len(tools)))
	return tools, nil
}

// DeleteTool 删除工具
func (s *ToolService) DeleteTool(ctx context.Context, toolID primitive.ObjectID) error {
	logger.Info(ctx, "Deleting tool started",
		slog.String("tool_id", toolID.Hex()))

	collection := database.GetCollection("tools")

	// 执行删除
	logger.Debug(ctx, "Executing tool deletion",
		slog.String("tool_id", toolID.Hex()))
	result, err := collection.DeleteOne(ctx, bson.M{"_id": toolID})
	if err != nil {
		logger.Error(ctx, "Failed to delete tool",
			slog.String("tool_id", toolID.Hex()),
			slog.Any("error", err))
		return err
	}

	// 检查是否删除了记录
	if result.DeletedCount == 0 {
		logger.Error(ctx, "Tool not found for deletion",
			slog.String("tool_id", toolID.Hex()))
		return ErrToolNotFound
	}

	logger.Info(ctx, "Tool deleted successfully",
		slog.String("tool_id", toolID.Hex()))
	return nil
}
