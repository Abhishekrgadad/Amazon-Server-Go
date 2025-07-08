package order

import (
	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(router fiber.Router) {

	order := router.Group("/order")
	order.Post("/checkout", CheckoutHandler)
	order.Get("/view", ViewOrdersHandler)
	order.Post("/cancel", CancelOrderHandler)
	order.Post("/return", ReturnOrderHandler)
	order.Post("/status", CheckOrderStatusHandler)
}
