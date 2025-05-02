package product

import (
	validation "server/modules/Validation"

	"github.com/gofiber/fiber/v2"
)

func AddProductHandler(c *fiber.Ctx) error {
	var product Product

	if err := c.BodyParser(&product); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"failed to parse json"})
	}

	if err := validation.ValidateInputs(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}
	
	err := AddProduct(&product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"failed to add product",})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message":"Product added successfully", "Product":product})
}

func ViewProductHandler(c *fiber.Ctx) error {

	products,err := ViewProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"failed to fetch the products"})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}