package models

import (
	"time"
)

// ChatMessage represents a chat message in the system
type ChatMessage struct {
	ID              string           `json:"id" bson:"_id,omitempty"`
	UserID          string           `json:"user_id" bson:"user_id"`
	Message         string           `json:"message" bson:"message"`
	Role            string           `json:"role" bson:"role"` // user or assistant
	Timestamp       time.Time        `json:"timestamp" bson:"timestamp"`
	EmotionAnalysis *EmotionAnalysis `json:"emotion_analysis,omitempty" bson:"emotion_analysis,omitempty"`
}

// EmotionAnalysis represents emotion analysis result
type EmotionAnalysis struct {
	Emotion    string  `json:"emotion" bson:"emotion"`
	Confidence float64 `json:"confidence" bson:"confidence"`
}
