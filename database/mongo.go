package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

const DBNAME = "webMessenger"
const DBNAME_HISTORY = "personalHistory"

// Устанавливает соединение с БД.
func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Проверка соединения.
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	MongoClient = client
	fmt.Println("Connected to MongoDB!")
}

// GetCollection возвращает коллекцию из базы данных DBNAME.
func GetCollection(collectionName string) *mongo.Collection {
	return MongoClient.Database(DBNAME).Collection(collectionName)
}

// GetCollectionHistory возвращает коллекцию из базы данных DBNAME_HISTORY
func GetCollectionHistory(collectionName string) *mongo.Collection {
	return MongoClient.Database(DBNAME_HISTORY).Collection(collectionName)
}
