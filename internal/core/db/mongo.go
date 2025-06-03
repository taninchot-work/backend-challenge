package db

import (
	"context"
	"fmt"
	"github.com/taninchot-work/backend-challenge/internal/core/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"slices"
	"time"
)

var database *mongo.Database

func InitializeMongoDB() error {
	mongoConfig := config.GetConfig().Database

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoConfig.ConnectionTimeout)*time.Millisecond)
	defer cancel()

	// validate configuration
	if mongoConfig.Host == "" || mongoConfig.Port == 0 || mongoConfig.DatabaseName == "" {
		return fmt.Errorf("invalid MongoDB configuration: %v", mongoConfig)
	}

	// create MongoDB client options
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", mongoConfig.Host, mongoConfig.Port)).
		SetMaxPoolSize(uint64(mongoConfig.MaxPoolSize))

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	database = client.Database(mongoConfig.DatabaseName)

	err = createCollection(ctx)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	err = createEmailUniqueIndex(ctx)
	if err != nil {
		return fmt.Errorf("failed to create unique index: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	return nil
}

func GetDatabase() *mongo.Database {
	if database == nil {
		log.Fatal("MongoDB is not initialized. Call InitializeMongoDB first.")
	}
	return database
}

func createEmailUniqueIndex(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"email", 1}},
		Options: options.Index().
			SetUnique(true).SetName("unique_email_index"),
	}

	indexView := database.Collection("users").Indexes()

	name, err := indexView.CreateOne(ctx, indexModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("Warning: Unique index on 'email' might already exist or conflicting data. Error: %v", err)
		} else {
			return fmt.Errorf("failed to create unique index on 'email' for 'users' collection: %w", err)
		}
	} else {
		log.Printf("Unique index '%s' created on 'email' for 'users' collection.", name)
	}
	return nil
}

func createCollection(ctx context.Context) error {
	// create collection if it doesn't exist
	collections, err := database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}
	if !slices.Contains(collections, "users") { // for one collection
		if err := database.CreateCollection(ctx, "users"); err != nil {
			return fmt.Errorf("failed to create 'users' collection: %v", err)
		}
		log.Println("Created 'users' collection in MongoDB")
	} else {
		log.Println("'users' collection already exists in MongoDB")
	}
	return nil
}
