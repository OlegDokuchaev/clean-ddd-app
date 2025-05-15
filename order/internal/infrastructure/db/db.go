package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewConnection(cfg *Config) (*mongo.Client, error) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(cfg.URI).SetConnectTimeout(cfg.ConnectTimeoutSeconds * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

func NewDB(cfg *Config, client *mongo.Client) *mongo.Database {
	return client.Database(cfg.Database)
}

func NewOrderCollection(db *mongo.Database, cfg *Config) *mongo.Collection {
	return db.Collection(cfg.OrderCollection)
}
