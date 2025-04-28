package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func ConnectDB() {

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("mongo uri is required")
	}

	clientoptions := options.Client().ApplyURI(uri)

	Client, err := mongo.Connect(context.TODO(), clientoptions)
	if err != nil {
		log.Fatal("Failed to connect to mongodb")
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("mongodb ping failed", err)
	}

	Database = Client.Database("Amazon")
	fmt.Println("MongoDB Connected Successfully")

}
