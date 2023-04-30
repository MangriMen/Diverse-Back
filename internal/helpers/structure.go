package helpers

import (
	"errors"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// GetNotEmpty returns the non-empty value between a and b.
func GetNotEmpty[T comparable](a T, b T) T {
	// newDeref lint error
	if a == *new(T) {
		return b
	}
	return a
}

// Ptr returns pointer to an object.
func Ptr[T any](obj T) *T { return &obj }

// CloseQuietly close io.Closer object quietly without returning error.
func CloseQuietly[T io.Closer](entity T) {
	if err := entity.Close(); err != nil {
		_ = entity.Close()
	}
}

// NewValidator creates a new instance of the validator.Validate struct,
// which is used to validate struct fields.
func NewValidator() *validator.Validate {
	validate := validator.New()

	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

// ValidatorErrors takes a validation error and returns a map of strings,
// where the keys represent the fields that failed validation
// and the values are the corresponding error messages.
func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	validationErrors := validator.ValidationErrors{}
	if !errors.As(err, &validationErrors) {
		return fields
	}

	for _, err := range validationErrors {
		fields[err.Field()] = err.Error()
	}

	return fields
}
