package router

import (
	auth "server/modules/Auth"
	product "server/modules/Product"

	"github.com/gofiber/fiber/v2"
)


func SetupRoutes(app *fiber.App) {
	root := app.Group("/api")
	auth.AuthRoutes(root)
	product.ProductRoutes(root)
}