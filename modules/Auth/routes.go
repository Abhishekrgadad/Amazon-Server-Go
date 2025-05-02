package auth

import "github.com/gofiber/fiber/v2"


func AuthRoutes(router fiber.Router){
	
	root := router.Group("/auth")
	root.Post("/register",RegisterHandler)
	root.Post("/login",LoginHandler)
	root.Get("/users",GetUserHandler)
	root.Get("/admins",GetAdminHandler)
}

