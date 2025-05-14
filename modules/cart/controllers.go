package cart

import (
	"server/config"
	"server/errors"
	"server/modules/product"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Function to add products to cart
func AddToCartHandler(c *fiber.Ctx) error {
	var req AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequestError(c, "Invalid request payload")
	}
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return errors.BadRequestError(c, "Invalid User ID")
	}

	var items []CartItem
	for _, item := range req.Items {
		pid, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return errors.BadRequestError(c, "Invalid Product ID")
		}
		items = append(items, CartItem{ProductID: pid, Quantity: item.Quantity})
	}
	cartInfo, err := AddToCart(userID, items)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"cart_id":     cartInfo.CartID.Hex(),
		"user_id":     userID,
		"items":       cartInfo.Details,
		"total_price": cartInfo.Total,
		"message":     "Added to cart successfully",
	})
}

// Function to display the cart details by id
func GetCartHandler(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "User ID is required"})
	}
	cart, err := GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(cart)
}

// Function to update the cart details
func UpdateCartHandler(c *fiber.Ctx) error {
	var req UpdateCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(struct {
			Error string `json:"error"`
		}{Error: "Invalid request payload"})
	}
	err := UpdateCartItem(req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	cart, err := GetCart(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct {
			Error string `json:"error"`
		}{Error: "Failed to update the cart"})
	}
	return c.Status(fiber.StatusOK).JSON(cart)
}

// Function to Remove Products from cart
func RemoveCartHandler(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	productID := c.Query("product_id")
	err := RemoveCartItem(userID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Item removed from cart"})
}

// Function to Remove Whole Cart from DB
func ClearCartHandler(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	err := ClearCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Cart cleared"})
}

// Updated GetAllCartsHandler to include detailed product information
func GetAllCartsHandler(c *fiber.Ctx) error {
	carts, err := GetAllCarts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var detailedCarts []fiber.Map
	for _, cart := range carts {
		var responseItems []CartItemResponse
		var total float64
		for _, item := range cart.Items {
			var product product.Product
			err := config.DB.Collection("products").FindOne(c.Context(), bson.M{"_id": item.ProductID}).Decode(&product)
			if err != nil {
				continue 
			}
			subtotal := product.Price * float64(item.Quantity)
			total += subtotal
			responseItems = append(responseItems, CartItemResponse{
				ProductID:   product.ID.Hex(),
				ProductName: product.Name,
				Price:       product.Price,
				Quantity:    item.Quantity,
				Subtotal:    subtotal,
			})
		}
		detailedCarts = append(detailedCarts, fiber.Map{
			"user_id":     cart.UserID.Hex(),
			"cart_id":     cart.ID.Hex(),
			"items":       responseItems,
			"total_price": total,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"carts": detailedCarts})
}
