package product

import (
	"context"
	"encoding/json"
	"fmt"
	"server/errors"
	"strconv"
	"time"

	"server/config"

	"github.com/gofiber/fiber/v2"
)

var Ctx = context.Background()

// Function to add a new product
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
	cacheKey := "product:id:" + result.InsertedID.(string)
	productJSON, err := json.Marshal(product)
	if err == nil {
		config.RedisClient.Set(Ctx, cacheKey, productJSON, time.Minute*10)
	}
	return c.JSON(fiber.Map{
		"message":    "Product added successfully",
		"product_id": result.InsertedID,
	})
}

// Function to get all products
func GetProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	cacheKey := "products:page:" + strconv.Itoa(page)
	cachedProducts, err := config.RedisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var productsFromCache []Product
		err = json.Unmarshal([]byte(cachedProducts), &productsFromCache)
		if err == nil {
			return c.JSON(fiber.Map{
				"data":        productsFromCache,
				"source":      "cache",
				"total_count": len(productsFromCache),
			})
		}
	}
	products, totalCount, totalPages, err := GetProducts(page)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	productsJSON, err := json.Marshal(products)
	if err != nil {
		fmt.Printf("Error serializing products for Redis: %v\n", err)
		return errors.InternalServerError(c, "Failed to cache products")
	}
	err = config.RedisClient.Set(Ctx, cacheKey, productsJSON, time.Minute*10).Err()
	if err != nil {
		config.ReconnectRedis()
		return errors.InternalServerError(c, "Failed to cache products")
	}
	return c.JSON(fiber.Map{
		"data":         products,
		"total_count":  totalCount,
		"total_pages":  totalPages,
		"current_page": page,
	})
}

// Function to get a product by ID
func GetProductByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	cacheKey := "product:id:" + id
	cachedProduct, err := config.RedisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var productFromCache Product
		if json.Unmarshal([]byte(cachedProduct), &productFromCache) == nil {
			return c.JSON(productFromCache)
		}
	}
	product, err := GetProductByID(id)
	if err != nil {
		return errors.NotFoundError(c, err.Error())
	}
	productJSON, err := json.Marshal(product)
	if err == nil {
		config.RedisClient.Set(Ctx, cacheKey, productJSON, time.Minute*10)
	}
	return c.JSON(product)
}

// Function to update a product
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
	cacheKey := "product:id:" + id
	config.RedisClient.Del(Ctx, cacheKey)
	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

// Function to delete a product
func DeleteProductHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := DeleteProduct(id)
	if err != nil {
		return errors.InternalServerError(c, err.Error())
	}
	cacheKey := "product:id:" + id
	config.RedisClient.Del(Ctx, cacheKey)
	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

// Function to get active products
func GetActiveProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	cacheKey := "active_products:page:" + strconv.Itoa(page)
	cachedProducts, err := config.RedisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var productsFromCache []Product
		err = json.Unmarshal([]byte(cachedProducts), &productsFromCache)
		if err == nil {
			return c.JSON(fiber.Map{
				"data":         productsFromCache,
				"source":       "cache",
				"current_page": page,
				"total_count":  len(productsFromCache),
			})
		}
	}
	products, totalCount, limit, err := GetActiveProducts(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch active products")
	}
	if len(products) == 0 {
		return errors.NotFoundError(c, "No active products found")
	}
	productsJSON, err := json.Marshal(products)
	if err == nil {
		config.RedisClient.Set(Ctx, cacheKey, productsJSON, time.Minute*10)
	}
	return c.JSON(fiber.Map{
		"products":    products,
		"message":     "Active products",
		"status":      "success",
		"limit":       limit,
		"total_count": totalCount,
	})
}

// Function to get inactive products
func GetInActiveProductsHandler(c *fiber.Ctx) error {
	pageStr := c.Params("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	cacheKey := "inactive_products:page:" + strconv.Itoa(page)
	cachedProducts, err := config.RedisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var productsFromCache []Product
		err = json.Unmarshal([]byte(cachedProducts), &productsFromCache)
		if err == nil {
			return c.JSON(fiber.Map{
				"data":         productsFromCache,
				"source":       "cache",
				"current_page": page,
				"total_count":  len(productsFromCache),
			})
		}
	}
	products, totalCount, err := GetInActiveProducts(page)
	if err != nil {
		return errors.InternalServerError(c, "Failed to fetch active products")
	}

	if len(products) == 0 {
		return errors.NotFoundError(c, "No active products found")
	}

	productsJSON, err := json.Marshal(products)
	if err == nil {
		config.RedisClient.Set(Ctx, cacheKey, productsJSON, time.Minute*10)
	}

	return c.JSON(fiber.Map{
		"products":    products,
		"message":     "inactive products",
		"status":      "success",
		"total_count": totalCount,
	})
}

// Function to apply different filter products
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
	cacheKey := fmt.Sprintf("filtered_products:name=%s:category=%s:brand=%s:minPrice=%.2f:maxPrice=%.2f", name, category, brand, minPrice, maxPrice)
	cachedProducts, err := config.RedisClient.Get(Ctx, cacheKey).Result()
	if err == nil {
		var productsFromCache []Product
		if json.Unmarshal([]byte(cachedProducts), &productsFromCache) == nil {
			return c.JSON(fiber.Map{"products": productsFromCache, "source": "cache"})
		}
	}
	products, err := FilteredProducts(name, category, brand, minPrice, maxPrice)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	productsJSON, err := json.Marshal(products)
	if err == nil {
		config.RedisClient.Set(Ctx, cacheKey, productsJSON, time.Minute*10)
	}
	return c.JSON(fiber.Map{"products": products})
}
