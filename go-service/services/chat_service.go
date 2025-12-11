package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"neuro-guide-go-service/database"
	"neuro-guide-go-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatService handles chat-related business logic
type ChatService struct {
	collection         *mongo.Collection
	pythonAIServiceURL string
	httpClient         *http.Client
}

// NewChatService creates a new instance of ChatService
func NewChatService(pythonAIServiceURL string) *ChatService {
	return &ChatService{
		collection:         database.Database.Collection("chat_messages"),
		pythonAIServiceURL: pythonAIServiceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatRequest represents a chat request to Python AI service
type ChatRequest struct {
	UserID  string                   `json:"user_id"`
	Message string                   `json:"message"`
	Context []map[string]interface{} `json:"context,omitempty"`
}

// ChatResponse represents a response from Python AI service
type ChatResponse struct {
	Response string `json:"response"`
}

// SendMessage sends a message to the Python AI service and saves it
func (cs *ChatService) SendMessage(userID, message string, context []map[string]interface{}) (string, error) {
	// Prepare request to Python AI service
	reqBody := ChatRequest{
		UserID:  userID,
		Message: message,
		Context: context,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Python AI service
	resp, err := cs.httpClient.Post(
		fmt.Sprintf("%s/chat", cs.pythonAIServiceURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI service returned error: %s", string(body))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Save user message
	userMsg := models.ChatMessage{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		Message:   message,
		Role:      "user",
		Timestamp: time.Now(),
	}
	if err := cs.SaveMessage(&userMsg); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save user message: %v\n", err)
	}

	// Save assistant response
	assistantMsg := models.ChatMessage{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		Message:   chatResp.Response,
		Role:      "assistant",
		Timestamp: time.Now(),
	}
	if err := cs.SaveMessage(&assistantMsg); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save assistant message: %v\n", err)
	}

	return chatResp.Response, nil
}

// SaveMessage saves a chat message to the database
func (cs *ChatService) SaveMessage(message *models.ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert string ID to ObjectID if needed
	if message.ID != "" {
		objID, err := primitive.ObjectIDFromHex(message.ID)
		if err == nil {
			_, err = cs.collection.InsertOne(ctx, bson.M{
				"_id":       objID,
				"user_id":   message.UserID,
				"message":   message.Message,
				"role":      message.Role,
				"timestamp": message.Timestamp,
			})
			return err
		}
	}

	_, err := cs.collection.InsertOne(ctx, message)
	return err
}

// GetChatHistory retrieves chat history for a user
func (cs *ChatService) GetChatHistory(userID string, limit int64) ([]*models.ChatMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(limit)

	cursor, err := cs.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.ChatMessage
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ClearChatHistory clears chat history for a user
func (cs *ChatService) ClearChatHistory(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	_, err := cs.collection.DeleteMany(ctx, filter)
	return err
}
