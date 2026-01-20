package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/models"
	"github.com/nspas/go-service/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIController struct {
	cfg                 *config.Config
	aiClient            services.AIClient
	conversationService *services.ConversationService
}

func NewAIController(cfg *config.Config) *AIController {
	return &AIController{
		cfg:                 cfg,
		aiClient:            services.NewHTTPClient(cfg),
		conversationService: services.NewConversationService(),
	}
}

// AIChatRequest AI聊天请求
type AIChatRequest struct {
	ConversationID   string             `json:"conversation_id,omitempty"`
	Messages         []Message          `json:"messages" binding:"required"`
	ConversationType string             `json:"conversation_type" binding:"required,oneof=analysis mapping assistant"`
}

// Message 消息结构
type Message struct {
	Content string `json:"content" binding:"required"`
	IsUser  bool   `json:"is_user" binding:"required"`
}

// AIChatResponse AI聊天响应
type AIChatResponse struct {
	Content string `json:"content"`
}

// Chat 处理AI对话请求
func (c *AIController) Chat(ctx *gin.Context) {
	var req AIChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		userIDStr = "anonymous" // 临时处理，后续需要添加认证
	}

	// 转换消息格式
	var aiMessages []services.Message
	for _, msg := range req.Messages {
		aiMessages = append(aiMessages, services.Message{
			Content: msg.Content,
			IsUser:  msg.IsUser,
		})
	}

	// 调用AI服务
	aiResponse, err := c.aiClient.Chat(ctx, aiMessages, req.ConversationType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call AI service"})
		return
	}

	// 创建完整的消息列表，包括AI响应
	var fullMessages []models.Message
	for _, msg := range req.Messages {
		fullMessages = append(fullMessages, models.Message{
			Content:   msg.Content,
			IsUser:    msg.IsUser,
			CreatedAt: time.Now(),
		})
	}

	// 添加AI响应
	fullMessages = append(fullMessages, models.Message{
		Content:   aiResponse,
		IsUser:    false,
		CreatedAt: time.Now(),
	})

	// 处理对话保存
	var conversation *models.Conversation
	if req.ConversationID == "" {
		// 创建新对话
		userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))
		conversation, err = c.conversationService.CreateConversation(ctx, userID, models.ConversationType(req.ConversationType), fullMessages[0].Content[:30])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
	} else {
		// 更新现有对话
		convID, err := primitive.ObjectIDFromHex(req.ConversationID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
			return
		}
		conversation, err = c.conversationService.GetConversationByID(ctx, convID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}
	}

	// 更新对话消息
	conversation.Messages = fullMessages
	conversation.UpdatedAt = time.Now()
	updatedConversation, err := c.conversationService.UpdateConversation(ctx, conversation.ID, fullMessages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
		return
	}

	// 返回响应
	ctx.JSON(http.StatusOK, gin.H{
		"content":         aiResponse,
		"conversation_id": updatedConversation.ID.Hex(),
		"messages":        fullMessages,
	})
}
