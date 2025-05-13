package review

import (
	"context"
	"fmt"
	"server/config"
	"server/errors"
	"server/modules/order"
	"server/modules/product"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddReview(req ReviewRequest) error {
    orderID, err := primitive.ObjectIDFromHex(req.OrderID)
    if err != nil {
        return err
    }
    productID, err := primitive.ObjectIDFromHex(req.ProductID)
    if err != nil {
        return err
    }
    orderCollection := config.DB.Collection("orders")
    var order order.Order
    err = orderCollection.FindOne(context.TODO(), bson.M{
        "_id": orderID,
        "items.product_id": productID,
    }).Decode(&order)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return errors.NotFoundError(nil, "Order or product not found")
        }
        return err
    }
    productFound := false
    for _, item := range order.Items {
        if item.ProductID == productID {
            productFound = true
            break
        }
    }
    if !productFound {
        return errors.NotFoundError(nil, "Product not found in the order")
    }

    review := Review{
        ID:        primitive.NewObjectID(),
        ProductID: productID,
        OrderID:   orderID,
        Rating:    req.Rating,
        Comment:   req.Comment,
        CreatedAt: time.Now(),
    }
    _, err = config.DB.Collection("reviews").InsertOne(context.TODO(), review)
    if err != nil {
        return err
    }
    productCollection := config.DB.Collection("products")
    var product product.Product
    err = productCollection.FindOne(context.TODO(), bson.M{"_id": productID}).Decode(&product)
    if err != nil {
        return err
    }
    newTotalReviews := product.TotalReviews + 1
    newAverageRating := ((product.AverageRating * float64(product.TotalReviews)) + float64(req.Rating)) / float64(newTotalReviews)
    update := bson.M{
        "$set": bson.M{
            "average_rating": newAverageRating,
            "total_reviews":  newTotalReviews,
        },
        "$push": bson.M{
            "comments": bson.M{
                "comment":    req.Comment,
                "rating":     req.Rating,
                "created_at": time.Now(),
            },
        },
    }
    _, err = productCollection.UpdateOne(context.TODO(), bson.M{"_id": productID}, update)
    if err != nil {
        return fmt.Errorf("failed to update product: %v", err)
    }
    return nil
}


func GetReviewsByProduct(productID primitive.ObjectID) ([]Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.DB.Collection("reviews").Find(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []Review
	if err = cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}
