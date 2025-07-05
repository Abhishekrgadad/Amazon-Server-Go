package product

import (
	"context"
	"server/errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var Ctx = context.Background()

func AddProductHandler(c *fiber.Ctx) error {
	var product Product
	if err := c.BodyParser(&product); err != nil {
		return errors.BadRequestError(c, "Invalid product details")
	}
	if err := ValidateProduct(&product); err != nil {
		return errors.BadRequestError(c, err.Error())
	}
	result, err := AddProduct(&product)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	return c.JSON(fiber.Map{
		"message":    "Product added successfully",
		"product_id": result.InsertedID,
	})
}

func GetProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	products, totalCount, totalPages, err := GetProducts(page)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"data":         products,
		"total_count":  totalCount,
		"total_pages":  totalPages,
		"current_page": page,
	})
}

func GetProductByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := GetProductByID(id)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	return c.JSON(product)
}

func UpdateProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var product Product
	if err := c.BodyParser(&product); err != nil {
		return errors.BadRequestError(c, "Invalid product data")
	}
	if err := ValidateProduct(&product); err != nil {
		return errors.BadRequestError(c, err.Error())
	}
	_, err := UpdateProduct(id, &product)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func DeleteProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := DeleteProduct(id)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

func GetActiveProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	products, totalCount, err := GetActiveProducts(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch products")
	}
	if len(products) == 0 {
		return errors.NotFoundError(c, "No active products found")
	}
	return c.JSON(fiber.Map{
		"products":    products,
		"message":     "Active products",
		"status":      "success",
		"total_count": totalCount,
	})
}

func GetInActiveProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	products, totalCount, err := GetInActiveProducts(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch products")
	}

	if len(products) == 0 {
		return errors.NotFoundError(c, "No active products found")
	}

	return c.JSON(fiber.Map{
		"products":    products,
		"message":     "inactive products",
		"status":      "success",
		"total_count": totalCount,
	})
}

func FilterProductsHandler(c *fiber.Ctx) error {
	name := c.Query("name")
	category := c.Query("category")
	brand := c.Query("brand")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	var minPrice, maxPrice float64
	var err error
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid min_price value"})
		}
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid max_price value"})
		}
	}

	products, err := FilteredProducts(name, category, brand, minPrice, maxPrice)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"products": products})
}
