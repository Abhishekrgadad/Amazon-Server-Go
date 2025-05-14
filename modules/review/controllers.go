package review

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddReviewHandler(c *fiber.Ctx) error {
	var request ReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	err := AddReview(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Review added successfully for the product"})
}

func GetProductReviewsHandler(c *fiber.Ctx) error {
	prodid := c.Params("product_id")
	productID, err := primitive.ObjectIDFromHex(prodid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	reviews, err := GetReviewsByProduct(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to fetch reviews"})
	}
	return c.JSON(reviews)
}

func UpdateReviewHandler(c *fiber.Ctx) error {
	var updateData ReviewUpdate
	err := c.BodyParser(&updateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse updatedata"})
	}
	productID, err := primitive.ObjectIDFromHex(updateData.ProductID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product id"})
	}
	UpdateReview, err := UpdateReview(productID, updateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to update review"})
	}
	return c.Status(fiber.StatusOK).JSON(UpdateReview)
}

func DeleteReviewHandler(c *fiber.Ctx) error {
	productID, err := primitive.ObjectIDFromHex(c.Params("product_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse the productid"})
	}
	err = DeleteReview(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to delete Review"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Review Deleted Successfully"})
}
