package product

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(router fiber.Router) {
	log.Println("Setting up product routes...")
	products := router.Group("/products") // Group all product-related routes under "/auth/products"
	products.Get("/true/page:page", GetActiveProductsHandler)
	products.Get("/false/page:page", GetInActiveProductsHandler)
	products.Get("/page:page", GetProductsHandler)       // Get all products with pagination
	products.Post("/add", AddProductHandler)             // Add new product
	products.Get("/:id", GetProductByIDHandler)          // Get product by ID
	products.Put("/update/:id", UpdateProductHandler)    // Update product by ID
	products.Delete("/delete/:id", DeleteProductHandler) // Delete product by ID
	log.Println("Registering /products/filter route...")
	products.Get("/filter", FilterProductsHandler)
}
