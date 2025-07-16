package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RequireRoles returns middleware that allows only specified roles
func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		role := claims["role"].(string)

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next() 
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Access denied. Contact adminf for more information.",
		})
	}
}
