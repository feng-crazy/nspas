package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tool struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name           string             `bson:"name" json:"name" validate:"required"`
	Description    string             `bson:"description" json:"description" validate:"required"`
	HTMLContent    string             `bson:"html_content" json:"html_content" validate:"required"`
	ConversationID primitive.ObjectID `bson:"conversation_id" json:"conversation_id"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}
