package errors

import "github.com/gofiber/fiber/v2"


func BadRequestError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
}

func UnauthorizedError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": message})
}

func ConflictError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": message})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": message})
}

func ForbiddenError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error":message})
}

func InvalidRoleError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":message})
}

func NotFoundError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":message})
}

