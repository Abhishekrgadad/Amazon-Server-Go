package product

import (
	validation "server/modules/Validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddProductHandler(c *fiber.Ctx) error {
	var product Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse json"})
	}

	if err := validation.ValidateInputs(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := AddProduct(&product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add product"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product added successfully", "Product": product})
}

func ViewProductHandler(c *fiber.Ctx) error {

	products, err := ViewProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch the products"})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

func FilterProductHandler(c *fiber.Ctx) error {
	category := c.Query("category")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	products, err := FilterProduct(category, minPrice, maxPrice)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to fetch products"})
	}
	return c.Status(fiber.StatusCreated).JSON(products)
}

func UpdateProductHandler(c *fiber.Ctx) error {

	id := c.Query("id")
	var updateData map[string]interface{}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse json"})
	}

	err := UpdateProduct(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to update products"})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func DeleteProductHandler(c *fiber.Ctx) error {

	id := c.Query("id")

	err := DeleteProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to delete product"})
	}
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}
