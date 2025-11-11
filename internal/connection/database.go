package connection

import (
	"context"
	"log"
	"time"

	"github.com/LouisFernando1204/kai-backend.git/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDatabase(conf config.Database) {
	clientOptions := options.Client().ApplyURI(conf.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error while connecting to MongoDB: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error while ping to MongoDB: %s", err)
	}

	Client = client
	log.Printf("Successfully connected to MongoDB!")
}

func GetCollection(dbName, collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatalf("Mongo DB is not initialized.")
	}
	return Client.Database(dbName).Collection(collectionName)
}

func GetClient() *mongo.Client {
	return Client
}

func DisconnectDatabase() {
	if Client == nil {
		log.Fatalf("MongoDB isn't initialized.")
	}

	err := Client.Disconnect(context.Background())
	if err != nil {
		log.Fatalf("Error while disconnecting to MongoDB: %s", err)
	}

	log.Printf("Successfully disconnected from MongoDB!")
}
