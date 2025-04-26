package review

import (
	"context"
	"server/config"
	"server/modules/order"
	"server/modules/product"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddReviewHandler(c *fiber.Ctx) error {
	var req ReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	orderCollection := config.DB.Collection("orders")
	var order order.Order
	err = orderCollection.FindOne(context.TODO(), bson.M{
		"_id":     orderID,
		"product_id": productID,
	}).Decode(&order)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Order not found or unauthorized"})
	}
	productFound := false
	for _, item := range order.Items {
		if item.ProductID == productID {
			productFound = true
			break
		}
	}
	if !productFound {
		return c.Status(400).JSON(fiber.Map{"error": "Product not found in the order"})
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
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add review"})
	}
	productCollection := config.DB.Collection("products")

	var product product.Product
	err = productCollection.FindOne(context.TODO(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}
	newTotalReviews := product.TotalReviews + 1
	newAverageRating := ((product.AverageRating * float64(product.TotalReviews)) + float64(req.Rating)) / float64(newTotalReviews)
	_, err = productCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": productID},
		bson.M{
			"$set": bson.M{
				"average_rating": newAverageRating,
				"total_reviews":  newTotalReviews,
			},
		},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update product rating"})
	}
	return c.JSON(fiber.Map{"message": "Review added successfully for the product"})
}

func GetProductReviewsHandler(c *fiber.Ctx) error {
	productID, err := primitive.ObjectIDFromHex(c.Params("product_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	reviews, err := GetReviewsByProduct(productID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reviews"})
	}
	return c.JSON(reviews)
}
