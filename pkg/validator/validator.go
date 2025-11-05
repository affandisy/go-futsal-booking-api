package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func New() *validator.Validate {
	validate = validator.New()
	return validate
}

func Validate(data interface{}) []string {
	var errors []string

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, formatError(err))
		}
	}

	return errors
}

func formatError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be after %s", field, err.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, err.Param())
	case "contains":
		return fmt.Sprintf("%s must contain '%s'", field, err.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain '%s'", field, err.Param())
	case "startswith":
		return fmt.Sprintf("%s must start with '%s'", field, err.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with '%s'", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
