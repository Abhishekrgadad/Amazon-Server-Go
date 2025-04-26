package product

import (
	"server/errors"
	"strconv"
	"time"
	"context"

	"server/config"

	"github.com/gofiber/fiber/v2"
)

var Ctx = context.Background()
// Function to add a product
func AddProductHandler(c *fiber.Ctx) error {
	var product Product
	if err := c.BodyParser(&product); err != nil {
		return errors.BadRequestError(c, "Invalid product details")
	}

	// Validate Product Data
	if err := ValidateProduct(&product); err != nil {
		return errors.BadRequestError(c, err.Error())
	}

	// Call function to add product
	result, err := AddProduct(&product)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"message":    "Product added successfully",
		"product_id": result.InsertedID,
	})
}


// Function to view products
func GetProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page") // Get the page parameter from URL

	// Convert page parameter to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1 if page parameter is invalid
	}

	// Check Redis cache for products
	cacheKey := "products:page:" + strconv.Itoa(page)
	cachedProducts, err := config.RedisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		// Return cached products
		return c.JSON(fiber.Map{
			"data":         cachedProducts,
			"source":       "cache",
			"current_page": page,
		})
	}

	// Fetch Products with Pagination
	products, totalCount, totalPages, err := GetProducts(page)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}

	// Cache the products in Redis
	err = config.RedisClient.Set(Ctx, cacheKey, products, time.Minute*10).Err()
	if err != nil {
		return errors.InternalServerError(c, "Failed to cache products")
	}

	// Return Products and Pagination Info
	return c.JSON(fiber.Map{
		"data":         products,
		"total_count":  totalCount,
		"total_pages":  totalPages,
		"current_page": page,
	})
}

// Get Product by ID Handler
func GetProductByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	// Fetch Product by ID
	product, err := GetProductByID(id)
	if err != nil {
		return errors.NotFoundError(c, err.Error())
	}

	return c.JSON(product)
}

// Update Product Handler
func UpdateProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var product Product
	if err := c.BodyParser(&product); err != nil {
		return errors.BadRequestError(c, "Invalid product data")
	}

	// Validate Product Data
	if err := ValidateProduct(&product); err != nil {
		return errors.BadRequestError(c, err.Error())
	}

	// Update Product
	_, err := UpdateProduct(id, &product)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

// Delete Product Handler
func DeleteProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	// Delete Product
	_, err := DeleteProduct(id)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// Active Products Handler
func GetActiveProductsHandler(c *fiber.Ctx) error {
	// Get active products from the database
	products, err := GetActiveProducts()
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch active products")
	}

	// Check if products are available
	if len(products) == 0 {
		return errors.NotFoundError(c, "No active products found")
	}

	return c.JSON(fiber.Map{
		"products": products,
		"message":  "Active products",
	})
}

// Inactive Products Handler
func GetInActiveProductsHandler(c *fiber.Ctx) error {
	// Get active products from the database
	products, err := GetInActiveProducts()
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch active products")
	}

	// Check if products are available
	if len(products) == 0 {
		return errors.NotFoundError(c, "No active products found")
	}

	return c.JSON(fiber.Map{
		"products": products,
		"message":  "inactive products",
	})
}

// FilterProductsHandler handles the API request for filtering products
func FilterProductsHandler(c *fiber.Ctx) error {
	name := c.Query("name")
	category := c.Query("category")
	brand := c.Query("brand")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	// Convert price values to float
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
	// Fetch products with filters
	products, err := FilteredProducts(name, category, brand, minPrice, maxPrice)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"products": products})
}
