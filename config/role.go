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
			"message": "Access denied. Contact admin for more information.",
		})
	}
}

// func RequireRoles(allowedRoles ...string) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Safely get token from context
// 		token := c.Locals("user")
// 		if token == nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Unauthorized. No token found.",
// 			})
// 		}

// 		// token, ok := tokenVal.(*jwt.Token)
// 		// if !ok {
// 		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 		// 		"success": false,
// 		// 		"message": "Invalid token format.",
// 		// 	})
// 		// }

// 		claims, ok := token.claims.(jwt.MapClaims)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Invalid claims format.",
// 			})
// 		}

// 		roleClaim, ok := claims["role"]
// 		if !ok {
// 			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Role claim missing in token.",
// 			})
// 		}

// 		role, ok := roleClaim.(string)
// 		if !ok {
// 			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 				"success": false,
// 				"message": "Invalid role claim.",
// 			})
// 		}

// 		for _, allowed := range allowedRoles {
// 			if role == allowed {
// 				return c.Next()
// 			}
// 		}

// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 			"success": false,
// 			"message": "Access denied. Contact admin for more information.",
// 		})
// 	}
// }
