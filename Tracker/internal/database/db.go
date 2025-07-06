package database

import (
	"context"
	"log"
	"time"

	"Tracker/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)

// InitDB initializes the database connection
func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB configuration
	mongoURI := config.GetMongoURI()
	dbName := config.GetMongoDBName()
	collectionName := config.GetMongoCollectionName()

	// Configure client options
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Printf("Connected to MongoDB at %s", mongoURI)

	// Initialize database and collection
	database = client.Database(dbName)
	collection = database.Collection(collectionName)

	return nil
}

// GetCollection returns the activities collection
func GetCollection() *mongo.Collection {
	if collection == nil {
		log.Fatal("Database collection not initialized. Call InitDB() first.")
	}
	return collection
}

// GetDatabase returns the database instance
func GetDatabase() *mongo.Database {
	if database == nil {
		log.Fatal("Database not initialized. Call InitDB() first.")
	}
	return database
}

// CloseDB closes the database connection
func CloseDB() error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("MongoDB connection closed")
	return nil
}