package controllers

import (
	"encoding/json"
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
	UserID           string    `json:"user_id" binding:"required"`
	ConversationID   string    `json:"conversation_id" binding:"required"`
	Messages         []Message `json:"messages" binding:"required"`
	ConversationType string    `json:"conversation_type" binding:"required,oneof=analysis mapping assistant"`
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
		slog.String("usr_id", req.UserID),
		slog.String("conversation_id", req.ConversationID),
		slog.String("conversation_type", req.ConversationType),
		slog.Int("message_count", len(req.Messages)))

	// 获取认证用户ID，如果不存在则返回错误
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		logger.Warn(reqCtx, "User not authenticated")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 确保请求的用户ID与认证的用户ID一致
	requestUserID := req.UserID
	authUserID := userIDStr.(string)
	if requestUserID != authUserID {
		logger.Warn(reqCtx, "User ID mismatch", 
			slog.String("request_user_id", requestUserID), 
			slog.String("authenticated_user_id", authUserID))
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User ID mismatch"})
		return
	}

	// 转换消息格式
	var aiMessages []services.Message
	for _, msg := range req.Messages {
		aiMessages = append(aiMessages, services.Message{
			Content: msg.Content,
			IsUser:  msg.IsUser,
		})
	}

	// 处理对话保存，获取或创建对话
	var conversation *models.Conversation
	var err error
	var isNewConversation bool

	if req.ConversationID == "" || req.ConversationID == "new" {
		// 创建新对话
		logger.Info(reqCtx, "Creating new conversation")
		userID, err := primitive.ObjectIDFromHex(requestUserID)
		if err != nil {
			logger.Warn(reqCtx, "Invalid user ID format", slog.String("user_id", requestUserID))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		
		title := req.Messages[0].Content
		if len(title) > 30 {
			title = title[:30] + "..."
		}
		conversation, err = c.conversationService.CreateConversation(reqCtx, userID, models.ConversationType(req.ConversationType), title)
		if err != nil {
			logger.Error(reqCtx, "Failed to create conversation", slog.Any("error", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
		isNewConversation = true
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
		
		// 验证对话属于当前用户
		if conversation.UserID.Hex() != requestUserID {
			logger.Warn(reqCtx, "User does not own conversation", 
				slog.String("request_user_id", requestUserID),
				slog.String("conversation_owner_id", conversation.UserID.Hex()))
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You don't own this conversation"})
			return
		}
	}

	// 创建完整的消息列表，包括AI响应占位符
	var fullMessages []models.Message
	for _, msg := range req.Messages {
		fullMessages = append(fullMessages, models.Message{
			Content:   msg.Content,
			IsUser:    msg.IsUser,
			CreatedAt: time.Now(),
		})
	}

	// 添加AI响应占位符（后续会实时更新）
	aiMessage := models.Message{
		Content:   "",
		IsUser:    false,
		CreatedAt: time.Now(),
	}
	fullMessages = append(fullMessages, aiMessage)

	// 设置SSE响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	// 调用AI服务的流式接口
	logger.Info(reqCtx, "Calling AI stream service")
	streamChan, err := c.aiClient.StreamChat(reqCtx, requestUserID, req.ConversationID, aiMessages, req.ConversationType)
	if err != nil {
		logger.Error(reqCtx, "Failed to call AI stream service", slog.Any("error", err))
		sseErr := map[string]any{
			"error": "Failed to call AI service",
		}
		if data, err := json.Marshal(sseErr); err == nil {
			ctx.SSEvent("error", string(data))
			ctx.Writer.Flush()
		}
		return
	}

	// 实时处理并发送AI响应
	var aiResponse string
	for chunk := range streamChan {
		aiResponse += chunk

		// 更新AI响应内容
		fullMessages[len(fullMessages)-1].Content = aiResponse

		// 构建SSE响应数据
		sseData := map[string]any{
			"content":           chunk,
			"full_content":      aiResponse,
			"conversation_id":   conversation.ID.Hex(),
			"is_new_conversation": isNewConversation,
			"messages":          fullMessages,
		}

		// 发送SSE事件
		if data, err := json.Marshal(sseData); err == nil {
			ctx.SSEvent("message", string(data))
			ctx.Writer.Flush()
		}
	}

	// 更新对话消息
	logger.Debug(reqCtx, "Updating conversation messages", slog.String("conversation_id", conversation.ID.Hex()))
	conversation.Messages = fullMessages
	conversation.UpdatedAt = time.Now()
	_, err = c.conversationService.UpdateConversation(reqCtx, conversation.ID, fullMessages)
	if err != nil {
		logger.Error(reqCtx, "Failed to update conversation", slog.Any("error", err))
	}

	// 发送完成事件
	sseComplete := map[string]any{
		"content":           aiResponse,
		"conversation_id":   conversation.ID.Hex(),
		"is_new_conversation": isNewConversation,
		"messages":          fullMessages,
		"completed":         true,
	}
	if data, err := json.Marshal(sseComplete); err == nil {
		ctx.SSEvent("complete", string(data))
		ctx.Writer.Flush()
	}

	logger.Info(reqCtx, "AI chat request completed successfully", slog.String("conversation_id", conversation.ID.Hex()))
}