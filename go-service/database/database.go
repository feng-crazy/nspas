package database

import (
	"context"
	"log"
	"time"

	"neuro-guide-go-service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database holds the mongoDB instance
var Database *mongo.Database

// InitDB initializes the database connection
func InitDB(cfg *config.Config) error {
	// MongoDB connection URI
	mongoURI := cfg.MongoDBURI
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // 默认MongoDB地址
	}

	// Database name
	dbName := cfg.MongoDBName
	if dbName == "" {
		dbName = "neuro_guide" // 默认数据库名
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")

	// Get a reference to the database
	Database = client.Database(dbName)

	return nil
}
