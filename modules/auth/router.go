package auth

import (
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes - Defines auth routes
func SetupAuthRoutes(router fiber.Router) {

	// Customer Routes
	users := router.Group("/users")
	users.Get("/page:page", GetAllUsersHandler)
	users.Get("/:id", GetUserHandler)
	users.Put("/update/:id", UpdateUserHandler)
	users.Delete("/delete/:id", DeleteUserHandler)

	// Seller Routes
	sellers := router.Group("/sellers")
	sellers.Get("/page:page", GetAllSellerHandler)
	sellers.Get("/:id", GetSellerHandler)
	sellers.Put("/update/:id", UpdateSellerHandler)
	sellers.Delete("/delete/:id", DeleteSellerHandler)

	// Admin Routes
	admins := router.Group("/admins")
	admins.Get("/page:page", GetAllAdminsHandler)
	admins.Get("/:id", GetAdminHandler)
	admins.Put("/update/:id", UpdateAdminHandler)
	admins.Delete("/delete/:id", DeleteAdminHandler)

	// Register Authentication Routes
	register := router.Group("/register")
	register.Post("/user", RegisterUserHandler)
	register.Post("/seller", RegisterSellerHandler)
	register.Post("/admin", RegisterAdminHandler)

	// Login Authentication Routes
	router.Post("/login", LoginHandler)

	// Password Reset Routes
	router.Post("/reset-password", RequestPasswordReset)
	router.Post("/update-password", UpdatePassword)
}
