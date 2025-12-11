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

// PracticeRecordService handles practice record-related business logic
type PracticeRecordService struct {
	collection *mongo.Collection
}

// NewPracticeRecordService creates a new instance of PracticeRecordService
func NewPracticeRecordService() *PracticeRecordService {
	return &PracticeRecordService{
		collection: database.Database.Collection("practice_records"),
	}
}

// CreateRecord creates a new practice record
func (prs *PracticeRecordService) CreateRecord(record *models.PracticeRecord) error {
	record.ID = primitive.NewObjectID().Hex()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(record.ID)
	if err != nil {
		return err
	}

	_, err = prs.collection.InsertOne(ctx, bson.M{
		"_id":             objID,
		"user_id":         record.UserID,
		"plan_id":         record.PlanID,
		"date":            record.Date,
		"completed_tasks": record.CompletedTasks,
		"reflection":      record.Reflection,
		"created_at":      record.CreatedAt,
		"updated_at":      record.UpdatedAt,
	})

	return err
}

// UpdateRecord updates a practice record
func (prs *PracticeRecordService) UpdateRecord(record *models.PracticeRecord) error {
	record.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(record.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"completed_tasks": record.CompletedTasks,
			"reflection":      record.Reflection,
			"updated_at":      record.UpdatedAt,
		},
	}

	_, err = prs.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// GetRecordByID retrieves a practice record by ID
func (prs *PracticeRecordService) GetRecordByID(id string) (*models.PracticeRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var record models.PracticeRecord
	err = prs.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// GetRecordsByUserID retrieves all practice records for a user
func (prs *PracticeRecordService) GetRecordsByUserID(userID string) ([]*models.PracticeRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.M{"date": -1})

	cursor, err := prs.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*models.PracticeRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecordsByPlanID retrieves all practice records for a plan
func (prs *PracticeRecordService) GetRecordsByPlanID(planID string) ([]*models.PracticeRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"plan_id": planID}
	opts := options.Find().SetSort(bson.M{"date": -1})

	cursor, err := prs.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*models.PracticeRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}
