package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/logger"
	"github.com/nspas/go-service/models"
	"github.com/nspas/go-service/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemoryController struct {
	conversationService *services.ConversationService
}

func NewMemoryController() *MemoryController {
	return &MemoryController{
		conversationService: services.NewConversationService(),
	}
}

// GetConversationMemory 获取对话记忆（最近3-5轮对话历史+长期记忆概要）
func (c *MemoryController) GetConversationMemory(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "Get conversation memory request started")

	// 获取对话ID
	conversationIDStr := ctx.Param("conversation_id")
	if conversationIDStr == "" {
		logger.Warn(reqCtx, "Missing conversation ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing conversation ID"})
		return
	}

	// 解析对话ID
	conversationID, err := primitive.ObjectIDFromHex(conversationIDStr)
	if err != nil {
		logger.Warn(reqCtx, "Invalid conversation ID", slog.String("conversation_id", conversationIDStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// 获取用户ID
	userIDStr := ctx.GetHeader("User-ID")
	if userIDStr == "" {
		logger.Warn(reqCtx, "Missing User-ID header")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing User-ID header"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		logger.Warn(reqCtx, "Invalid User-ID", slog.String("user_id", userIDStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User-ID"})
		return
	}

	// 获取对话
	conversation, err := c.conversationService.GetConversationByID(reqCtx, conversationID)
	if err != nil {
		logger.Error(reqCtx, "Failed to get conversation", slog.Any("error", err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	// 验证对话属于该用户
	if conversation.UserID != userID {
		logger.Warn(reqCtx, "Conversation does not belong to user",
			slog.String("conversation_id", conversationIDStr),
			slog.String("user_id", userIDStr))
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Conversation does not belong to user"})
		return
	}

	// 准备记忆数据
	memory := struct {
		RecentMessages  []models.Message `json:"recent_messages"`
		LongTermSummary string           `json:"long_term_summary"`
	}{
		RecentMessages:  []models.Message{},
		LongTermSummary: "", // 目前长期记忆概要为空，后续可以从数据库中获取
	}

	// 获取最近3-5轮对话
	messageCount := len(conversation.Messages)
	startIndex := 0
	if messageCount > 10 { // 如果对话超过10轮，只返回最近5轮
		startIndex = messageCount - 10
	} else if messageCount > 6 { // 如果对话超过6轮，返回最近5轮
		startIndex = messageCount - 10
	} else if messageCount > 3 { // 如果对话超过3轮，返回最近3轮
		startIndex = messageCount - 6
	}

	if startIndex < 0 {
		startIndex = 0
	}

	// 添加最近消息
	memory.RecentMessages = conversation.Messages[startIndex:]

	logger.Info(reqCtx, "Get conversation memory request completed successfully",
		slog.String("conversation_id", conversationIDStr),
		slog.Int("recent_messages_count", len(memory.RecentMessages)))

	ctx.JSON(http.StatusOK, memory)
}
