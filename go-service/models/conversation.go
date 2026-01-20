package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationType string

const (
	ConversationTypeAnalysis  ConversationType = "analysis"
	ConversationTypeMapping   ConversationType = "mapping"
	ConversationTypeAssistant ConversationType = "assistant"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Content   string             `bson:"content" json:"content"`
	IsUser    bool               `bson:"is_user" json:"is_user"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Type      ConversationType   `bson:"type" json:"type"`
	Title     string             `bson:"title" json:"title"`
	Messages  []Message          `bson:"messages" json:"messages"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
