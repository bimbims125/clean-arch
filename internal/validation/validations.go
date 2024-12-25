package validation

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string][]string {
	errors := make(map[string][]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := fieldErr.Field()
			message := getErrorMessage(fieldErr)
			errors[field] = append(errors[field], message)
		}
	}
	return errors
}

// ValidatePassword checks if a password has at least one uppercase letter and one special character
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := false
	hasSpecial := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if strings.ContainsRune("!@#$%^&*()_+-={}[]|:;\"'<>,.?/~`", char) {
			hasSpecial = true
		} else if unicode.IsNumber(char) {
			hasNumber = true
		}
		if hasUpper && hasSpecial && hasNumber {
			return true
		}
	}

	return false
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "This field must be a valid email address"
	case "alphanum":
		return "This field must be alphanumeric"
	case "min":
		return "This field must be at least " + fe.Param() + " characters"
	case "password":
		return "This field must contain at least one uppercase letter, one number, and one special character"
	case "unique":
		return "Email already exists"
	default:
		return "Invalid value"
	}
}
