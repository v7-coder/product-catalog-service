package validator

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		validate.SetTagName("validate")
	})

	return validate
}

func Validate(s interface{}) error {
	err := getValidator().Struct(s)
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range errs {
			messages = append(messages, formatError(e))
		}
		return fmt.Errorf("%s", strings.Join(messages, "; "))
	}

	return err
}

func formatError(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	default:
		return fmt.Sprintf("%s validation failed: %s=%s", field, e.Tag(), e.Param())
	}
}
