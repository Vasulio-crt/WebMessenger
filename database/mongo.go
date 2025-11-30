package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectDB устанавливает соединение с MongoDB.
func ConnectDB() {
	// Установите URI для вашего MongoDB.
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

// GetCollection возвращает коллекцию из базы данных.
func GetCollection(dbName, collectionName string) *mongo.Collection {
	return MongoClient.Database(dbName).Collection(collectionName)
}

// CreateIndexes создает индексы для коллекций в MongoDB.
func CreateIndexes() {
	messagesCollection := GetCollection("webMessenger", "messages")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "timestamp", Value: -1}}, // -1 для сортировки по убыванию
		Options: options.Index().SetName("timestamp_desc"),
	}

	_, err := messagesCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Println("Failed to create index for messages:", err)
	}
	fmt.Println("Created index for messages collection.")
}
