package database

import (
	"context"
	"log"
	"time"

	"github.com/nspas/go-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	log.Println("Connected to MongoDB")
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("nspas").Collection(collectionName)
}

func Close() error {
	if Client != nil {
		return Client.Disconnect(context.Background())
	}
	return nil
}
