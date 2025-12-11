package models

import (
	"time"
)

// PracticeRecord represents a practice record in the system
type PracticeRecord struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	UserID         string    `json:"user_id" bson:"user_id"`
	PlanID         string    `json:"plan_id" bson:"plan_id"`
	Date           time.Time `json:"date" bson:"date"`
	CompletedTasks []string  `json:"completed_tasks" bson:"completed_tasks"`
	Reflection     string    `json:"reflection" bson:"reflection"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}
