package controllers

import (
	"net/http"

	"neuro-guide-go-service/config"
	"neuro-guide-go-service/services"

	"github.com/gin-gonic/gin"
)

var chatService *services.ChatService

// InitChatController initializes the chat controller with config
func InitChatController(cfg *config.Config) {
	chatService = services.NewChatService(cfg.PythonAIServiceURL)
}

// ChatMessageRequest represents a chat message request
type ChatMessageRequest struct {
	Message string                   `json:"message" binding:"required"`
	Context []map[string]interface{} `json:"context,omitempty"`
}

// SendMessage handles sending a chat message
func SendMessage(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send message to Python AI service without adding chat history context
	// The Python service will manage context internally
	response, err := chatService.SendMessage(userID, req.Message, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

// GetChatHistory handles getting chat history
func GetChatHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limit := int64(50) // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		// Parse limit from query string if provided
		// For simplicity, using default
	}

	messages, err := chatService.GetChatHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// ClearChatHistory handles clearing chat history
func ClearChatHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := chatService.ClearChatHistory(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat history cleared"})
}
