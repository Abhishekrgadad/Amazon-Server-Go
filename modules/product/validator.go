package product

import ("errors"
	"fmt"
)


func ValidateProduct(product *Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Description == "" {
		return errors.New("product description is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.StockQuantity < 0 {
		return errors.New("product stock quantity cannot be negative")
	}
	if product.Category == "" {
		return errors.New("product category is required")
	}
	if product.Brand == "" {
		return errors.New("product category is required")
	}
	return nil
}

func ValidateProductVisibility(visibility string) error {
	validVisibilities := []string{"active", "inactive"}
	for _, v := range validVisibilities {
		if visibility == v {
			return nil
		}
	}
	return fmt.Errorf("invalid visibility status: must be 'active' or 'inactive'")
}
