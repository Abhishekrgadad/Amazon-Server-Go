package product

import (
	"context"
	"server/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProduct(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product.ID = primitive.NewObjectID()
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := config.DB.Collection("products")
	_, err := collection.InsertOne(ctx, product)
	return err
}

func ViewProduct() ([]Product, error) {
	collecion := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collecion.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func FilterProduct(category string, minPrice, maxPrice float64) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	priceFilter := bson.M{}
	if minPrice > 0 {
		priceFilter["minPrice"] = minPrice
	}
	if maxPrice > 0 {
		priceFilter["maxPrice"] = maxPrice
	}
	if len(priceFilter) > 0 {
		filter["price"] = priceFilter
	}

	collection := config.DB.Collection("products")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			continue
		}
		products = append(products, product)
	}
	return products, nil
}
