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

func CheckoutHandler(c *fiber.Ctx) error {
	var request CheckoutRequest
	if err := c.BodyParser(&request); err != nil {
		return errors.InternalServerError(c, "Invalid Request")
	}
	userID, _ := primitive.ObjectIDFromHex(request.UserID)
	cartID, _ := primitive.ObjectIDFromHex(request.CartID)
	order, err := PlaceOrder(userID, cartID, request.PaymentType, request.Address, request.CouponCode)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	OrderStatus(order.ID)
	deliveryDate := time.Now().Add(48 * time.Hour).Format("02-Jan-2006")
	websocket.SendOrderNotification(fmt.Sprintf("Order placed successfully!\n It will be delivered on time.\nOrder ID: %s,\nUser ID: %s", order.ID.Hex(), userID.Hex()))
	return c.JSON(fiber.Map{
		"message":           "Order Placed Successfully",
		"expected_delivery": deliveryDate,
		"order":             order,
	})
}

func ViewOrdersHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")

	userID, err := primitive.ObjectIDFromHex(c.Query("user_id"))
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

func CancelOrderHandler(c *fiber.Ctx) error {
	orderID, err := primitive.ObjectIDFromHex(c.Query("order_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid order ID"})
	}
	userID, err := primitive.ObjectIDFromHex(c.Query("user_id"))
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

func ReturnOrderHandler(c *fiber.Ctx) error {
	orderID, err := primitive.ObjectIDFromHex(c.Query("order_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	err = ReturnOrder(orderID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	websocket.CancelOrderNotification(fmt.Sprintf("Order Returned initiated now.\n Refund will be initiated after the product received.\nUser ID: %s", orderID.Hex()))
	return c.JSON(fiber.Map{"message": "Order returned successfully"})
}

func CheckOrderStatusHandler(c *fiber.Ctx) error {
	orderIDParam := c.Query("order_id")
	orderID, err := primitive.ObjectIDFromHex(orderIDParam)
	if err != nil {
		return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	
	var order Order
	collection := config.DB.Collection("orders")
	err = collection.FindOne(ctx,bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"order_id": orderID.Hex(),
		"status":   order.Status,
		"success": true,
	})
}
