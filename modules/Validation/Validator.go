package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init(){
	validate = validator.New()
}

func ValidateInputs(User interface{}) error {
	if err := validate.Struct(User); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func formatValidationError(err error) error {
	if _,ok := err.(validator.ValidationErrors); ok {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field '%s' failed validation on '%s' ", e.Field(), e.Tag())
		}
	}
	return err
}