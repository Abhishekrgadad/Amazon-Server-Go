package product

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"server/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProduct(product *Product) (*mongo.InsertOneResult, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingProduct Product
	err := collection.FindOne(ctx, bson.M{"name": product.Name}).Decode(&existingProduct)
	if err == nil {
		return nil, fmt.Errorf("product already exists")
	}
	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to add product")
	}
	return result, nil
}

func UpdateProduct(id string, product *Product) (*mongo.UpdateResult, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}
	product.UpdatedAt = time.Now().Format(time.RFC3339)
	updateResult, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": product},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}
	return updateResult, nil
}

func DeleteProduct(id string) (*mongo.DeleteResult, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %v", err)
	}
	return deleteResult, nil
}

func GetProducts(page int) ([]Product, int64, int, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := int64(10)
	skip := int64(page-1) * limit
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count products: %v", err)
	}
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch products: %v", err)
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to parse data: %v", err)
	}
	totalPages := int(totalCount / limit)
	if totalCount%limit != 0 {
		totalPages++
	}
	return products, totalCount, totalPages, nil
}

func GetProductByID(id string) (*Product, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}
	var product Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}
	return &product, nil
}

func GetActiveProducts(page int) ([]Product, int64, int64, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := int64(10)
	skip := int64(page-1) * limit
	filter := bson.M{"visibility": true}
	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count products: %v", err)
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch products: %v", err)
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to parse product data: %v", err)
	}
	return products, totalCount, limit, nil
}

func GetInActiveProducts(page int) ([]Product, int64, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := int64(10)
	skip := int64(page-1) * limit
	filter := bson.M{"visibility": false}
	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %v", err)
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{Key: "name",Value: 1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %v", err)
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("failed to parse product data: %v", err)
	}
	return products, totalCount, nil
}

func FilteredProducts(name, category, brand string, minPrice, maxPrice float64) ([]Product, error) {
	collection := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"visibility": true}
	if category != "" {
		filter["category"] = category
	}
	if brand != "" {
		filter["brand"] = brand
	}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if minPrice > 0 && maxPrice > 0 {
		filter["price"] = bson.M{"$gte": minPrice, "$lte": maxPrice}
	} else if minPrice > 0 {
		filter["price"] = bson.M{"$gte": minPrice}
	} else if maxPrice > 0 {
		filter["price"] = bson.M{"$lte": maxPrice}
	}
	limit := int64(10)
	opts := options.Find().SetSort(bson.D{{Key: "name",Value: 1}}).SetLimit(limit)
	cursor, err := collection.Find(ctx, filter,opts)
	if err != nil {
		log.Printf("MongoDB Find error: %v", err)
		return nil, fmt.Errorf("failed to query products")
	}
	defer cursor.Close(ctx)

	var products []Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("MongoDB Cursor decode error: %v", err)
		return nil, fmt.Errorf("failed to decode products data")
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("no products found matching the criteria")
	}
	return products, nil
}
