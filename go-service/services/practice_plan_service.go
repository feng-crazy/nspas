package services

import (
	"context"
	"time"

	"neuro-guide-go-service/database"
	"neuro-guide-go-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PracticePlanService handles practice plan-related business logic
type PracticePlanService struct {
	collection *mongo.Collection
}

// NewPracticePlanService creates a new instance of PracticePlanService
func NewPracticePlanService() *PracticePlanService {
	return &PracticePlanService{
		collection: database.Database.Collection("practice_plans"),
	}
}

// CreatePlan creates a new practice plan
func (pps *PracticePlanService) CreatePlan(plan *models.PracticePlan) error {
	plan.ID = primitive.NewObjectID().Hex()
	plan.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(plan.ID)
	if err != nil {
		return err
	}

	_, err = pps.collection.InsertOne(ctx, bson.M{
		"_id":        objID,
		"user_id":    plan.UserID,
		"title":      plan.Title,
		"days":       plan.Days,
		"tasks":      plan.Tasks,
		"created_at": plan.CreatedAt,
	})

	return err
}

// GetPlanByID retrieves a practice plan by ID
func (pps *PracticePlanService) GetPlanByID(id string) (*models.PracticePlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var plan models.PracticePlan
	err = pps.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

// GetPlansByUserID retrieves all practice plans for a user
func (pps *PracticePlanService) GetPlansByUserID(userID string) ([]*models.PracticePlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := pps.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []*models.PracticePlan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

// DeletePlan deletes a practice plan
func (pps *PracticePlanService) DeletePlan(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = pps.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
