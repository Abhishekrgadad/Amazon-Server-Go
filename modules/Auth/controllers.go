package auth

import (
	
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	
	if err := c.BodyParser(&RegisterRequest); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Invalid request body"})
	}

	switch RegisterRequest.Role {
	case "user":
		return RegisterUserHandler(c)
	case "admin":
		return RegisterAdmin(c)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Role"})
	}
}