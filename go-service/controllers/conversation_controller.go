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

type ConversationController struct {
	conversationService *services.ConversationService
}

func NewConversationController() *ConversationController {
	return &ConversationController{
		conversationService: services.NewConversationService(),
	}
}

// CreateConversationRequest 创建对话请求
type CreateConversationRequest struct {
	Type  models.ConversationType `json:"type" binding:"required,oneof=analysis mapping assistant"`
	Title string                  `json:"title" binding:"required"`
}

// UpdateConversationRequest 更新对话请求
type UpdateConversationRequest struct {
	Messages []models.Message `json:"messages" binding:"required"`
}

// CreateConversation 创建对话
func (c *ConversationController) CreateConversation(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "Create conversation request started")

	var req CreateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn(reqCtx, "Invalid create conversation request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		logger.Warn(reqCtx, "User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		logger.Warn(reqCtx, "Invalid user ID", slog.String("user_id", userIDStr.(string)))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	logger.Debug(reqCtx, "Creating conversation",
		slog.String("user_id", userID.Hex()),
		slog.String("type", string(req.Type)),
		slog.String("title", req.Title))

	// 创建对话
	conversation, err := c.conversationService.CreateConversation(reqCtx, userID, req.Type, req.Title)
	if err != nil {
		logger.Error(reqCtx, "Failed to create conversation", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	logger.Info(reqCtx, "Conversation created successfully", slog.String("conversation_id", conversation.ID.Hex()))
	ctx.JSON(http.StatusCreated, conversation)
}

// UpdateConversation 更新对话
func (c *ConversationController) UpdateConversation(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	convIDStr := ctx.Param("id")
	logger.Info(reqCtx, "Update conversation request started", slog.String("conversation_id", convIDStr))

	var req UpdateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warn(reqCtx, "Invalid update conversation request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取对话ID
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		logger.Warn(reqCtx, "Invalid conversation ID", slog.String("conversation_id", convIDStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	logger.Debug(reqCtx, "Updating conversation messages",
		slog.String("conversation_id", convID.Hex()),
		slog.Int("message_count", len(req.Messages)))

	// 更新对话
	conversation, err := c.conversationService.UpdateConversation(reqCtx, convID, req.Messages)
	if err != nil {
		if err == services.ErrConversationNotFound {
			logger.Warn(reqCtx, "Conversation not found", slog.String("conversation_id", convID.Hex()))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.Error(reqCtx, "Failed to update conversation", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
		return
	}

	logger.Info(reqCtx, "Conversation updated successfully", slog.String("conversation_id", conversation.ID.Hex()))
	ctx.JSON(http.StatusOK, conversation)
}

// GetConversation 获取对话
func (c *ConversationController) GetConversation(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	convIDStr := ctx.Param("id")
	logger.Info(reqCtx, "Get conversation request started", slog.String("conversation_id", convIDStr))

	// 获取对话ID
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		logger.Warn(reqCtx, "Invalid conversation ID", slog.String("conversation_id", convIDStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	logger.Debug(reqCtx, "Getting conversation by ID", slog.String("conversation_id", convID.Hex()))

	// 获取对话
	conversation, err := c.conversationService.GetConversationByID(reqCtx, convID)
	if err != nil {
		if err == services.ErrConversationNotFound {
			logger.Warn(reqCtx, "Conversation not found", slog.String("conversation_id", convID.Hex()))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.Error(reqCtx, "Failed to get conversation", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		return
	}

	logger.Info(reqCtx, "Got conversation successfully", slog.String("conversation_id", conversation.ID.Hex()))
	ctx.JSON(http.StatusOK, conversation)
}

// GetUserConversations 获取用户的所有对话
func (c *ConversationController) GetUserConversations(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	logger.Info(reqCtx, "Get user conversations request started")

	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		logger.Warn(reqCtx, "User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		logger.Warn(reqCtx, "Invalid user ID", slog.String("user_id", userIDStr.(string)))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	logger.Debug(reqCtx, "Getting user conversations", slog.String("user_id", userID.Hex()))

	// 获取用户对话
	conversations, err := c.conversationService.GetUserConversations(reqCtx, userID)
	if err != nil {
		logger.Error(reqCtx, "Failed to get user conversations", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversations"})
		return
	}

	logger.Info(reqCtx, "Got user conversations successfully",
		slog.String("user_id", userID.Hex()),
		slog.Int("conversation_count", len(conversations)))
	ctx.JSON(http.StatusOK, conversations)
}

// DeleteConversation 删除对话
func (c *ConversationController) DeleteConversation(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	convIDStr := ctx.Param("id")
	logger.Info(reqCtx, "Delete conversation request started", slog.String("conversation_id", convIDStr))

	// 获取对话ID
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		logger.Warn(reqCtx, "Invalid conversation ID", slog.String("conversation_id", convIDStr))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	logger.Debug(reqCtx, "Deleting conversation", slog.String("conversation_id", convID.Hex()))

	// 删除对话
	err = c.conversationService.DeleteConversation(reqCtx, convID)
	if err != nil {
		if err == services.ErrConversationNotFound {
			logger.Warn(reqCtx, "Conversation not found", slog.String("conversation_id", convID.Hex()))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.Error(reqCtx, "Failed to delete conversation", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		return
	}

	logger.Info(reqCtx, "Conversation deleted successfully", slog.String("conversation_id", convID.Hex()))
	ctx.JSON(http.StatusNoContent, nil)
}
