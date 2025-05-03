package product

import (
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(router fiber.Router) {
	root := router.Group("/products")
	root.Post("/add", AddProductHandler)
	root.Get("/view", ViewProductHandler)
	root.Get("/filter", FilterProductHandler)
	root.Put("/update",UpdateProductHandler)
	root.Delete("/delete",DeleteProductHandler)
}
