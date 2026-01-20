package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	var req CreateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 创建对话
	conversation, err := c.conversationService.CreateConversation(ctx, userID, req.Type, req.Title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	ctx.JSON(http.StatusCreated, conversation)
}

// UpdateConversation 更新对话
func (c *ConversationController) UpdateConversation(ctx *gin.Context) {
	var req UpdateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取对话ID
	convIDStr := ctx.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// 更新对话
	conversation, err := c.conversationService.UpdateConversation(ctx, convID, req.Messages)
	if err != nil {
		if err == services.ErrConversationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
		return
	}

	ctx.JSON(http.StatusOK, conversation)
}

// GetConversation 获取对话
func (c *ConversationController) GetConversation(ctx *gin.Context) {
	// 获取对话ID
	convIDStr := ctx.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// 获取对话
	conversation, err := c.conversationService.GetConversationByID(ctx, convID)
	if err != nil {
		if err == services.ErrConversationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		return
	}

	ctx.JSON(http.StatusOK, conversation)
}

// GetUserConversations 获取用户的所有对话
func (c *ConversationController) GetUserConversations(ctx *gin.Context) {
	// 从上下文获取用户ID
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 解析用户ID
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 获取用户对话
	conversations, err := c.conversationService.GetUserConversations(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversations"})
		return
	}

	ctx.JSON(http.StatusOK, conversations)
}

// DeleteConversation 删除对话
func (c *ConversationController) DeleteConversation(ctx *gin.Context) {
	// 获取对话ID
	convIDStr := ctx.Param("id")
	convID, err := primitive.ObjectIDFromHex(convIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// 删除对话
	err = c.conversationService.DeleteConversation(ctx, convID)
	if err != nil {
		if err == services.ErrConversationNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete conversation"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
