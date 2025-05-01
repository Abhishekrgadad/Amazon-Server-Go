package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
) 

var DB *mongo.Database

func ConnectDB(){

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	client,err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil{
		log.Fatal(err)
	}
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil{
		log.Fatal(err)
	}

	DB = client.Database(dbName)
	log.Println("Connected to Mongodb")
}