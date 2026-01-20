package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nspas/go-service/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToolController struct {
	toolService *services.ToolService
}

func NewToolController() *ToolController {
	return &ToolController{
		toolService: services.NewToolService(),
	}
}

// SaveToolRequest 保存工具请求
type SaveToolRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	HTMLContent    string `json:"html_content" binding:"required"`
	ConversationID string `json:"conversation_id" binding:"required"`
}

// SaveTool 保存工具
func (c *ToolController) SaveTool(ctx *gin.Context) {
	var req SaveToolRequest
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

	// 解析对话ID
	convID, err := primitive.ObjectIDFromHex(req.ConversationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	// 保存工具
	tool, err := c.toolService.SaveTool(ctx, userID, req.Name, req.Description, req.HTMLContent, convID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save tool"})
		return
	}

	ctx.JSON(http.StatusCreated, tool)
}

// GetTool 获取工具详情
func (c *ToolController) GetTool(ctx *gin.Context) {
	// 获取工具ID
	toolIDStr := ctx.Param("id")
	toolID, err := primitive.ObjectIDFromHex(toolIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tool ID"})
		return
	}

	// 获取工具
	tool, err := c.toolService.GetToolByID(ctx, toolID)
	if err != nil {
		if err == services.ErrToolNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tool"})
		return
	}

	ctx.JSON(http.StatusOK, tool)
}

// GetUserTools 获取用户的所有工具
func (c *ToolController) GetUserTools(ctx *gin.Context) {
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

	// 获取用户工具
	tools, err := c.toolService.GetUserTools(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tools"})
		return
	}

	ctx.JSON(http.StatusOK, tools)
}

// DeleteTool 删除工具
func (c *ToolController) DeleteTool(ctx *gin.Context) {
	// 获取工具ID
	toolIDStr := ctx.Param("id")
	toolID, err := primitive.ObjectIDFromHex(toolIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tool ID"})
		return
	}

	// 删除工具
	err = c.toolService.DeleteTool(ctx, toolID)
	if err != nil {
		if err == services.ErrToolNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tool"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
