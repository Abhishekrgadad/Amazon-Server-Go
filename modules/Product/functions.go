package product

import (
	"context"
	"server/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProduct(product *Product) error {
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	product.ID = primitive.NewObjectID()
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := config.DB.Collection("products")
	_,err := collection.InsertOne(ctx,product)
	return err
}

func ViewProduct() ([]Product,error) {
	collecion := config.DB.Collection("products")
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	cursor,err := collecion.Find(ctx,bson.M{})
	if err != nil {
		return nil,err
	}
	defer cursor.Close(ctx)

	var products []Product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil{
			continue
		}
		products = append(products, product)
	}
	return products,nil
}