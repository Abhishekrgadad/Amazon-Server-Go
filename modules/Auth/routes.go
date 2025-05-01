package auth

import "github.com/gofiber/fiber/v2"


func AuthRoutes(router fiber.Router){
	
	register := router.Group("/register")
	register.Post("/user",Register)
}

