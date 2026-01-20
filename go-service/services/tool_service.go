package services

import (
	"context"
	"time"

	"github.com/nspas/go-service/database"
	"github.com/nspas/go-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToolService struct {}

func NewToolService() *ToolService {
	return &ToolService{}
}

// SaveTool 保存工具
func (s *ToolService) SaveTool(ctx context.Context, userID primitive.ObjectID, name, description, htmlContent string, convID primitive.ObjectID) (*models.Tool, error) {
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
	_, err := collection.InsertOne(ctx, tool)
	if err != nil {
		return nil, err
	}

	return tool, nil
}

// GetToolByID 根据ID获取工具
func (s *ToolService) GetToolByID(ctx context.Context, toolID primitive.ObjectID) (*models.Tool, error) {
	collection := database.GetCollection("tools")

	var tool models.Tool
	err := collection.FindOne(ctx, bson.M{"_id": toolID}).Decode(&tool)
	if err != nil {
		return nil, ErrToolNotFound
	}

	return &tool, nil
}

// GetUserTools 获取用户的所有工具
func (s *ToolService) GetUserTools(ctx context.Context, userID primitive.ObjectID) ([]models.Tool, error) {
	collection := database.GetCollection("tools")

	// 查询用户的所有工具
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历结果
	var tools []models.Tool
	for cursor.Next(ctx) {
		var tool models.Tool
		if err := cursor.Decode(&tool); err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}

	return tools, nil
}

// DeleteTool 删除工具
func (s *ToolService) DeleteTool(ctx context.Context, toolID primitive.ObjectID) error {
	collection := database.GetCollection("tools")

	// 执行删除
	result, err := collection.DeleteOne(ctx, bson.M{"_id": toolID})
	if err != nil {
		return err
	}

	// 检查是否删除了记录
	if result.DeletedCount == 0 {
		return ErrToolNotFound
	}

	return nil
}
