package review

import "github.com/gofiber/fiber/v2"

func SetupReviewRoutes(router fiber.Router) {
	review := router.Group("/review")
	review.Post("/add", AddReviewHandler)
	review.Get("/view/:product_id", GetProductReviewsHandler)
	review.Put("/update", UpdateReviewHandler)
	review.Delete("/delete/:product_id",DeleteReviewHandler)
}
