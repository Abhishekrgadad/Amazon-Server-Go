package product

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(router fiber.Router) {
	products := router.Group("/products") 
	products.Get("/true/page:page", GetActiveProductsHandler)
	products.Get("/false/page:page", GetInActiveProductsHandler)
	products.Get("/page:page", GetProductsHandler)       
	products.Post("/add", AddProductHandler)            
	products.Get("/:id", GetProductByIDHandler)         
	products.Put("/update/:id", UpdateProductHandler)   
	products.Delete("/delete/:id", DeleteProductHandler) 
	fmt.Println("at filter route")
	products.Get("/filter", FilterProductsHandler)
}
