package coupons

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCouponHandler(c *fiber.Ctx) error {
	var coupon Coupon
	if err := c.BodyParser(&coupon); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
	}
	if err := CreateCoupon(&coupon); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create coupon"})
	}
	return c.JSON(coupon)
}

func GetAllCouponsHandler(c *fiber.Ctx) error {
	coupons, err := GetAllCoupons()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch coupons"})
	}
	return c.JSON(coupons)
}

func GetCouponByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid coupon ID"})
	}
	coupon, err := GetCouponByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Coupon not found"})
	}
	return c.JSON(coupon)
}

func UpdateCouponHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid coupon ID"})
	}

	var updates bson.M
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
	}
	if _, ok := updates["expiry_date"]; ok {
		if dateStr, ok := updates["expiry_date"].(string); ok {
			t, _ := time.Parse(time.RFC3339, dateStr)
			updates["expiry_date"] = t
		}
	}
	if err := UpdateCoupon(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update coupon"})
	}
	return c.JSON(fiber.Map{"message": "Coupon updated"})
}

func DeleteCouponHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid coupon ID"})
	}
	if err := DeleteCoupon(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete coupon"})
	}
	return c.JSON(fiber.Map{"message": "Coupon deleted"})
}
