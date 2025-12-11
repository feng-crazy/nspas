package models

import (
	"time"
)

// PracticePlan represents a practice plan in the system
type PracticePlan struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	UserID    string     `json:"user_id" bson:"user_id"`
	Title     string     `json:"title" bson:"title"`
	Days      int        `json:"days" bson:"days"`
	Tasks     []PlanTask `json:"tasks" bson:"tasks"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

// PlanTask represents a task in a practice plan
type PlanTask struct {
	Day             int    `json:"day" bson:"day"`
	Title           string `json:"title" bson:"title"`
	Description     string `json:"description" bson:"description"`
	ScientificBasis string `json:"scientific_basis" bson:"scientific_basis"`
}
