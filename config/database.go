package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB connection failed: ", err)
	}

	DB = client.Database("amazon_platform")

	fmt.Println("MongoDB Connected")
}
