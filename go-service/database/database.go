package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info(ctx, "Initializing MongoDB connection", slog.String("uri", cfg.Database.URI))

	clientOptions := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error(ctx, "Failed to create MongoDB client", slog.Any("error", err))
		return err
	}

	// 测试连接
	logger.Info(ctx, "Testing MongoDB connection")
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Error(ctx, "MongoDB ping failed", slog.Any("error", err))
		return err
	}

	Client = client
	logger.Info(ctx, "Connected to MongoDB successfully")
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	ctx := context.Background()
	logger.Debug(ctx, "Getting MongoDB collection", slog.String("collection", collectionName))
	return Client.Database("nspas").Collection(collectionName)
}

func Close() error {
	ctx := context.Background()
	if Client != nil {
		logger.Info(ctx, "Closing MongoDB connection")
		err := Client.Disconnect(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to close MongoDB connection", slog.Any("error", err))
			return err
		}
		logger.Info(ctx, "MongoDB connection closed successfully")
	}
	return nil
}
