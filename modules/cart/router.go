package cart

import (
	"github.com/gofiber/fiber/v2"
)

func SetupCartRoutes(router fiber.Router) {

	cart := router.Group("/cart")
	cart.Post("/add", AddToCartHandler)     
	cart.Get("/view", GetCartHandler)         
	cart.Put("/update", UpdateCartHandler)    
	cart.Delete("/delete", RemoveCartHandler) 
	cart.Delete("/clear", ClearCartHandler)  
}
