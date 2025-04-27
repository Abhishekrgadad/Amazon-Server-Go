package router

import (
	"server/modules/auth"
	"server/modules/cart"
	"server/modules/coupons"
	"server/modules/order"
	"server/modules/product"
	"server/modules/review"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	root := app.Group("/auth")
	auth.SetupAuthRoutes(root)
	product.SetupProductRoutes(root)
	cart.SetupCartRoutes(root)
	order.SetupOrderRoutes(root)
	review.SetupReviewRoutes(root)
	coupons.SetupCouponRoutes(root)
}
