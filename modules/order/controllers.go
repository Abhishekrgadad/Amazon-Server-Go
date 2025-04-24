package order

import (
	"context"
	"fmt"
	"server/config"
	"server/errors"
	"server/modules/websocket"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Function to place the order
func CheckoutHandler(c *fiber.Ctx) error {
	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.InternalServerError(c, "Invalid Request")
	}
	userID, _ := primitive.ObjectIDFromHex(req.UserID)
	cartID, _ := primitive.ObjectIDFromHex(req.CartID)

	order, err := PlaceOrder(userID, cartID, req.PaymentType, req.Address, req.CouponCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	deliveryDate := time.Now().Add(48 * time.Hour).Format("02-Jan-2025")

	websocket.SendOrderNotification(fmt.Sprintf("Order placed successfully!\n It will be delivered on time.\nOrder ID: %s", order.ID.Hex()))

	return c.JSON(fiber.Map{
		"message":           "Order Placed Successfully",
		"expected_delivery": deliveryDate,
		"order":             order,
	})
}

// Function for order status
func OrderStatusHandler(c *fiber.Ctx) error {
	orderID := c.Params("order_id")

	go func(orderID string) {
		time.Sleep(10 * time.Second)
		UpdateOrderStatus(orderID, "Shipped")
		time.Sleep(10 * time.Second)
		UpdateOrderStatus(orderID, "Out for Delivery")
		time.Sleep(10 * time.Second)
		UpdateOrderStatus(orderID, "Delivered")
	}(orderID)

	return c.JSON(fiber.Map{
		"status": "Order Confirmed",
		"date":   time.Now().Format("02-Jan-2006 15:04:05"),
	})
}

// Function for updating order status
func UpdateOrderStatus(orderID string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderObjID, _ := primitive.ObjectIDFromHex(orderID)
	_, err := config.DB.Collection("orders").UpdateOne(ctx, bson.M{"_id": orderObjID}, bson.M{"$set": bson.M{"status": status}})
	return err
}

// Function to view placed orders
func ViewOrdersHandler(c *fiber.Ctx) error {
	var req ViewOrdersRequest
	pageStr := c.Params("page")
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequestError(c, "Invalid request body")
	}
	if req.UserID == "" {
		return errors.BadRequestError(c, "user_id is required")
	}
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return errors.BadRequestError(c, "Invalid user_id")
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	orders, totalCount, totalPages, err := ViewAllOrders(userID, page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch orders")
	}
	if len(orders) == 0 {
		return errors.NotFoundError(c, "No orders found for this user")
	}
	return c.JSON(fiber.Map{
		"orders":       orders,
		"total_count":  totalCount,
		"total_pages":  totalPages,
		"current_page": page,
	})
}

// Function to cancel order
func CancelOrderHandler(c *fiber.Ctx) error {
	var req CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid order ID"})
	}
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user ID"})
	}
	err = CancelOrder(orderID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	websocket.CancelOrderNotification(fmt.Sprintf("Order Cancelled Successfully.\n Refund will be initiated soon.\nOrder ID: %s", orderID.Hex()))
	return c.JSON(fiber.Map{"message": "Order cancelled successfully. Refund will be initiated soon"})
}

// Function for return the order/products.
func ReturnOrderHandler(c *fiber.Ctx) error {
	var req ReturnRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	orderID, err := primitive.ObjectIDFromHex(req.OrderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	err = ReturnOrder(orderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	websocket.CancelOrderNotification(fmt.Sprintf("📦 Order Returned initiated now.\n Refund will be initiated after the product received.\nOrder ID: %s", orderID.Hex()))
	return c.JSON(fiber.Map{"message": "Order returned successfully"})
}
