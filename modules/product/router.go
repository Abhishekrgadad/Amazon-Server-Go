package product

import (
	"server/config"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(router fiber.Router) {
	products := router.Group("/products")
	products.Get("/true/page:page", GetActiveProductsHandler)
	products.Get("/false/page:page", GetInActiveProductsHandler)

	products.Get("/page:page", config.RequireRoles("admin","seller"),GetProductsHandler)
	products.Post("/add", AddProductHandler)
	products.Get("/get/:id", GetProductByIDHandler)
	products.Put("/update/:id", UpdateProductHandler)
	products.Delete("/delete/:id", DeleteProductHandler)
	products.Get("/filter", FilterProductsHandler)
}
