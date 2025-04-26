package review

import (
	"context"
	"fmt"
	"server/config"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddReview(userID, productID, orderID primitive.ObjectID, rating int, comment string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	review := Review{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID,
		OrderID:   orderID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	if _, err := config.DB.Collection("reviews").InsertOne(ctx, review); err != nil {
		return fmt.Errorf("failed to insert review: %v", err)
	}
	reviews, err := GetReviewsByProduct(productID)
	if err != nil {
		return fmt.Errorf("failed to fetch reviews: %v", err)
	}

	var totalRating int
	var comments []string
	for i := 0; i < len(reviews); i++ {
		totalRating += reviews[i].Rating
	}

	// Get latest 5 comments (sorted by CreatedAt descending)
	sort.SliceStable(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
	for i := 0; i < len(reviews) && i < 5; i++ {
		comments = append(comments, reviews[i].Comment)
	}
	avgRating := float64(totalRating) / float64(len(reviews))
	// Update the product document with average rating and latest comments
	update := bson.M{
		"$set": bson.M{
			"average_rating":  avgRating,
			"review_comments": comments,
		},
	}
	_, err = config.DB.Collection("products").UpdateOne(ctx, bson.M{"_id": productID}, update)
	if err != nil {
		return fmt.Errorf("failed to update product with review data: %v", err)
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
