package coupons

import "github.com/gofiber/fiber/v2"

func SetupCouponRoutes(router fiber.Router) {
	coupon := router.Group("/coupons")

	coupon.Post("/add", CreateCouponHandler)
	coupon.Get("/view", GetAllCouponsHandler)
	coupon.Get("/view/:id", GetCouponByIDHandler)
	coupon.Put("/update/:id", UpdateCouponHandler)
	coupon.Delete("/delete/:id", DeleteCouponHandler)
}
