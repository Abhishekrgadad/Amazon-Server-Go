package review

import "github.com/gofiber/fiber/v2"

func SetupReviewRoutes(router fiber.Router) {
	review := router.Group("/review")
	review.Post("/add", AddReviewHandler)
	review.Get("/:product_id", GetProductReviewsHandler)
}
