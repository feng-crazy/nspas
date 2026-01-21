package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
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
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "AI chat request started")

	var req AIChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn(reqCtx, "Invalid AI chat request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug(reqCtx, "Processing AI chat request",
		slog.String("conversation_id", req.ConversationID),
		slog.String("conversation_type", req.ConversationType),
		slog.Int("message_count", len(req.Messages)))

	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		userIDStr = "anonymous" // 临时处理，后续需要添加认证
		logger.Warn(reqCtx, "User not authenticated, using anonymous")
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
	logger.Info(reqCtx, "Calling AI service")
	aiResponse, err := c.aiClient.Chat(reqCtx, aiMessages, req.ConversationType)
	if err != nil {
		logger.Error(reqCtx, "Failed to call AI service", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call AI service"})
		return
	}
	logger.Debug(reqCtx, "AI service response received", slog.String("response", aiResponse[:50]+"..."))

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
		logger.Info(reqCtx, "Creating new conversation")
		userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))
		conversation, err = c.conversationService.CreateConversation(reqCtx, userID, models.ConversationType(req.ConversationType), fullMessages[0].Content[:30])
		if err != nil {
			logger.Error(reqCtx, "Failed to create conversation", slog.Any("error", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
		logger.Info(reqCtx, "Conversation created successfully", slog.String("conversation_id", conversation.ID.Hex()))
	} else {
		// 更新现有对话
		logger.Info(reqCtx, "Updating existing conversation", slog.String("conversation_id", req.ConversationID))
		convID, err := primitive.ObjectIDFromHex(req.ConversationID)
		if err != nil {
			logger.Warn(reqCtx, "Invalid conversation ID", slog.String("conversation_id", req.ConversationID))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
			return
		}
		conversation, err = c.conversationService.GetConversationByID(reqCtx, convID)
		if err != nil {
			logger.Error(reqCtx, "Conversation not found", slog.Any("error", err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}
	}

	// 更新对话消息
	logger.Debug(reqCtx, "Updating conversation messages", slog.String("conversation_id", conversation.ID.Hex()))
	conversation.Messages = fullMessages
	conversation.UpdatedAt = time.Now()
	updatedConversation, err := c.conversationService.UpdateConversation(reqCtx, conversation.ID, fullMessages)
	if err != nil {
		logger.Error(reqCtx, "Failed to update conversation", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
		return
	}

	// 返回响应
	logger.Info(reqCtx, "AI chat request completed successfully", slog.String("conversation_id", updatedConversation.ID.Hex()))
	ctx.JSON(http.StatusOK, gin.H{
		"content":         aiResponse,
		"conversation_id": updatedConversation.ID.Hex(),
		"messages":        fullMessages,
	})
}
