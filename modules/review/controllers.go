package review

import (

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddReviewHandler(c *fiber.Ctx) error {
	var req ReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	err := AddReview(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
