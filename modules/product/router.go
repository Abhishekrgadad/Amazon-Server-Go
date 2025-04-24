package product

import "github.com/gofiber/fiber/v2"

// 📌 Setup Product Routes
func SetupProductRoutes(router fiber.Router) {
	products := router.Group("/products") // Group all product-related routes under "/auth/products"
	products.Get("/true", GetActiveProductsHandler)
	products.Get("/false", GetInActiveProductsHandler)
	products.Get("/page:page", GetProductsHandler)       // Get all products with pagination
	products.Post("/add", AddProductHandler)             // Add new product
	products.Get("/:id", GetProductByIDHandler)          // Get product by ID
	products.Put("/update/:id", UpdateProductHandler)    // Update product by ID
	products.Delete("/delete/:id", DeleteProductHandler) // Delete product by ID
	products.Get("/filter", FilterProductsHandler)
}
