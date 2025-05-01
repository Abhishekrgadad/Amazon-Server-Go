package auth

import (

	validation "server/modules/Validation"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandler(c *fiber.Ctx) error {
	var RegisterRequest RegisterRequest
	if err := c.BodyParser(&RegisterRequest); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Invalid request body"})
	}

	switch RegisterRequest.Role {
	case "user":
		return RegisterUserHandler(c)
	case "admin":
		return RegisterAdminHandler(c)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Role"})
	}
}

func LoginHandler(c *fiber.Ctx) error {
	var input LoginRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Cannot parse JSON"})
	}
	if err := validation.ValidateInputs(&input); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	token,err := Login(input)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token": token,
	})


}