package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required."
	case "email":
		return "Invalid email format."
	case "min":
		minLength, _ := strconv.Atoi(fe.Param())
		return fmt.Sprintf("Value is too short. Minimum %d characters are required.", minLength)
	case "max":
		maxLength, _ := strconv.Atoi(fe.Param())
		return fmt.Sprintf("Value is too long. Maximum %d characters are allowed.", maxLength)
	case "alpha":
		return "Only alphabetic characters are allowed."
	case "alphanum":
		return "Only alphanumeric characters are allowed."
	case "boolean":
		return "Value must be a boolean."
	case "number":
		return "Value must be a number."
	case "numeric":
		return "Value must be a numeric."
	case "json":
		return "Value must be a valid JSON."
	case "file":
		return "Value must be a file."
	case "url":
		return "Value must be a valid URL."
	case "base64":
		return "Value must be a valid Base64 string."
	default:
		return fe.Error() // Fallback to the default error message if no custom message is specified
	}
}

func ValidationErrors(err error, c *gin.Context) []APIError {
	var ve validator.ValidationErrors
	var out []APIError
	if errors.As(err, &ve) {
		out = make([]APIError, len(ve))
		for i, fe := range ve {
			out[i] = APIError{fe.Field(), MsgForTag(fe)}
		}
	}
	return out
}
