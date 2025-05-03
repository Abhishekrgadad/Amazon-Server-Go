package auth

import "github.com/gofiber/fiber/v2"


func AuthRoutes(router fiber.Router){
	
	root := router.Group("/auth")
	root.Post("/register",RegisterHandler)
	root.Post("/login",LoginHandler)

	users := router.Group("/users")
	users.Get("/",GetUserHandler)
	users.Put("/update",UpdateUserHandler)
	users.Delete("/delete",DeleteUserHandler)

	admins := router.Group("/admins")
	admins.Get("/",GetAdminHandler)
	admins.Put("/update",UpdateAdminHandler)
	admins.Delete("/delete",DeleteAdminHandler)
}

