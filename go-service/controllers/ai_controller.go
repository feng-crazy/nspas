package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
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
	ConversationID   string    `json:"conversation_id"`
	Messages         []Message `json:"messages" binding:"required"`
	ConversationType string    `json:"conversation_type" binding:"required,oneof=analysis mapping assistant"`
}

// Message 消息结构
type Message struct {
	Content string `json:"content" binding:"required"`
	IsUser  bool   `json:"isUser" binding:"required"`
}

// Chat 处理AI对话请求
func (c *AIController) Chat(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "AI chat request started")

	// 解析和验证请求
	req, err := c.parseAndValidateRequest(ctx, reqCtx)
	if err != nil {
		return
	}

	// 验证用户认证
	requestUserID, err := c.validateUserAuthentication(ctx, reqCtx, req)
	if err != nil {
		return
	}

	// 转换消息格式
	aiMessages := c.convertMessages(req.Messages)

	// 处理对话（创建或更新）
	conversation, err := c.handleConversation(ctx, reqCtx, requestUserID, req)
	if err != nil {
		return
	}

	// 准备完整的消息列表
	fullMessages := c.prepareFullMessages(req.Messages, *conversation)

	// 设置SSE响应头
	c.setupSSEResponse(ctx)

	// 调用AI服务并处理响应
	aiResponse, err := c.processAIResponse(ctx, reqCtx, requestUserID, req, aiMessages)
	if err != nil {
		return
	}

	// 更新对话消息
	logger.Debug(reqCtx, "Updating conversation messages", slog.String("conversation_id", conversation.ID.Hex()))
	fullMessages[len(fullMessages)-1].Content = aiResponse
	conversation.Messages = fullMessages
	conversation.UpdatedAt = time.Now()
	_, err = c.conversationService.UpdateConversation(reqCtx, conversation.ID, fullMessages)
	if err != nil {
		logger.Error(reqCtx, "Failed to update conversation", slog.Any("error", err))
	}

	logger.Info(reqCtx, "AI chat request completed successfully", slog.String("conversation_id", conversation.ID.Hex()))
}

// parseAndValidateRequest 解析和验证请求
func (c *AIController) parseAndValidateRequest(ctx *gin.Context, reqCtx context.Context) (*AIChatRequest, error) {
	var req AIChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn(reqCtx, "Invalid AI chat request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	logger.Debug(reqCtx, "Processing AI chat request",
		slog.String("usr_id", req.UserID),
		slog.String("conversation_id", req.ConversationID),
		slog.String("conversation_type", req.ConversationType),
		slog.Int("message_count", len(req.Messages)))

	return &req, nil
}

// validateUserAuthentication 验证用户认证
func (c *AIController) validateUserAuthentication(ctx *gin.Context, reqCtx context.Context, req *AIChatRequest) (string, error) {
	// 获取认证用户ID，如果不存在则返回错误
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		logger.Warn(reqCtx, "User not authenticated")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return "", errors.New("user not authenticated")
	}

	// 确保请求的用户ID与认证的用户ID一致
	requestUserID := req.UserID
	authUserID := userIDStr.(string)
	if requestUserID != authUserID {
		logger.Warn(reqCtx, "User ID mismatch",
			slog.String("request_user_id", requestUserID),
			slog.String("authenticated_user_id", authUserID))
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User ID mismatch"})
		return "", errors.New("user ID mismatch")
	}

	return requestUserID, nil
}

// convertMessages 转换消息格式
func (c *AIController) convertMessages(messages []Message) []services.Message {
	var aiMessages []services.Message
	for _, msg := range messages {
		aiMessages = append(aiMessages, services.Message{
			Content: msg.Content,
			IsUser:  msg.IsUser,
		})
	}
	return aiMessages
}

// handleConversation 处理对话（创建或更新）
func (c *AIController) handleConversation(ctx *gin.Context, reqCtx context.Context,
	requestUserID string, req *AIChatRequest) (*models.Conversation, error) {

	var conversation *models.Conversation
	var err error
	var isNewConversation bool

	if req.ConversationID == "" || req.ConversationID == "new" {
		// 创建新对话
		conversation, isNewConversation, err = c.createNewConversation(ctx, reqCtx, requestUserID, req)
	} else {
		// 更新现有对话
		conversation, err = c.getExistingConversation(ctx, reqCtx, requestUserID, req)
	}

	if isNewConversation {
		req.ConversationID = conversation.ID.Hex()
	}

	return conversation, err
}

// createNewConversation 创建新对话
func (c *AIController) createNewConversation(ctx *gin.Context, reqCtx context.Context, requestUserID string, req *AIChatRequest) (*models.Conversation, bool, error) {
	logger.Info(reqCtx, "Creating new conversation")
	userID, err := primitive.ObjectIDFromHex(requestUserID)
	if err != nil {
		logger.Warn(reqCtx, "Invalid user ID format", slog.String("user_id", requestUserID))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return nil, false, err
	}

	title := req.Messages[0].Content
	if len(title) > 30 {
		title = title[:30] + "..."
	}
	conversation, err := c.conversationService.CreateConversation(reqCtx, userID, models.ConversationType(req.ConversationType), title)
	if err != nil {
		logger.Error(reqCtx, "Failed to create conversation", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return nil, false, err
	}
	isNewConversation := true
	logger.Info(reqCtx, "Conversation created successfully", slog.String("conversation_id", conversation.ID.Hex()))

	return conversation, isNewConversation, nil
}

// getExistingConversation 更新现有对话
func (c *AIController) getExistingConversation(ctx *gin.Context, reqCtx context.Context, requestUserID string, req *AIChatRequest) (*models.Conversation, error) {
	logger.Info(reqCtx, "Updating existing conversation", slog.String("conversation_id", req.ConversationID))
	convID, err := primitive.ObjectIDFromHex(req.ConversationID)
	if err != nil {
		logger.Warn(reqCtx, "Invalid conversation ID", slog.String("conversation_id", req.ConversationID))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return nil, err
	}

	conversation, err := c.conversationService.GetConversationByID(reqCtx, convID)
	if err != nil {
		logger.Error(reqCtx, "Conversation not found", slog.Any("error", err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return nil, err
	}

	// 验证对话属于当前用户
	if conversation.UserID.Hex() != requestUserID {
		logger.Warn(reqCtx, "User does not own conversation",
			slog.String("request_user_id", requestUserID),
			slog.String("conversation_owner_id", conversation.UserID.Hex()))
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You don't own this conversation"})
		return nil, err
	}

	return conversation, nil
}

// prepareFullMessages 准备完整的消息列表
func (c *AIController) prepareFullMessages(messages []Message, conversation models.Conversation) []models.Message {
	var fullMessages []models.Message
	for _, msg := range conversation.Messages {
		fullMessages = append(fullMessages, models.Message{
			Content:   msg.Content,
			IsUser:    msg.IsUser,
			CreatedAt: time.Now(),
		})
	}

	for _, msg := range messages {
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

	return fullMessages
}

// setupSSEResponse 设置SSE响应头
func (c *AIController) setupSSEResponse(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
}

// processAIResponse 调用AI服务并处理响应
func (c *AIController) processAIResponse(ctx *gin.Context, reqCtx context.Context, requestUserID string,
	req *AIChatRequest, aiMessages []services.Message) (string, error) {

	// 调用AI服务的流式接口
	logger.Info(reqCtx, "Calling AI stream service")
	streamChan, err := c.aiClient.StreamChat(reqCtx, requestUserID,
		req.ConversationID, aiMessages,
		req.ConversationType)
	if err != nil {
		logger.Error(reqCtx, "Failed to call AI stream service", slog.Any("error", err))
		c.sendSSEError(ctx, "Failed to call AI service")
		return "", err
	}

	// 实时处理并发送AI响应
	aiResponse, err := c.handleStreamResponse(ctx, streamChan)
	if err != nil {
		return "", err
	}

	return aiResponse, nil
}

// sendSSEError 发送SSE错误事件
func (c *AIController) sendSSEError(ctx *gin.Context, errorMsg string) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	// 直接发送错误信息，使用SSE格式
	ctx.Writer.Write([]byte("event: error\n"))
	ctx.Writer.Write([]byte("data: {\"error\": \"" + errorMsg + "\"}\n\n"))
	ctx.Writer.Flush()
}

// SSEData SSE响应数据结构
type SSEData struct {
	Content string `json:"content"`
}

// parseSSEData 解析SSE数据并返回content字段
func (c *AIController) parseSSEData(jsonData string) (string, error) {
	var sseData SSEData
	if err := json.Unmarshal([]byte(jsonData), &sseData); err != nil {
		return "", err
	}
	return sseData.Content, nil
}

// handleStreamResponse 处理流式响应
func (c *AIController) handleStreamResponse(ctx *gin.Context, streamChan <-chan string) (string, error) {
	var fullContent strings.Builder

	for chunk := range streamChan {
		// 直接透传数据，保留原始格式（包括data:前缀）
		ctx.Writer.Write([]byte(chunk + "\n\n"))
		ctx.Writer.Flush()

		// 提取实际内容用于保存到对话历史
		if strings.HasPrefix(chunk, "data: ") {
			jsonData := strings.TrimSpace(chunk[6:])

			content, err := c.parseSSEData(jsonData)
			if err == nil {
				// 成功解析JSON，追加content内容
				fullContent.WriteString(content)
			} else {
				// 如果不是标准JSON格式，检查是否是结束标记或其他类型的数据
				if jsonData != "[DONE]" && !strings.Contains(jsonData, "\"error\":") {
					logger.Warn(context.Background(), "Failed to parse SSE data JSON", slog.String("data", jsonData), slog.Any("error", err))
				}
			}
		}
	}

	return fullContent.String(), nil
}
