package helper

import (
	"bytes"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate = validator.New()

// ValidationError represents a single field error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResponse is a unified error format containing one or more errors.
type ValidationResponse struct {
	Errors []ValidationError `json:"errors"`
}

func ExtractValidateRequestBody(data any, ctx *fiber.Ctx) (bool, *ValidationResponse) {
	if err := parseRequestBodyStrict(data, ctx); err != nil {
		return false, &ValidationResponse{
			Errors: []ValidationError{
				{
					Field:   "Body parser",
					Message: err.Error(),
				},
			},
		}
	}

	if err := Validate.Struct(data); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var errors []ValidationError
			for _, fe := range ve {
				errors = append(errors, ValidationError{
					Field:   fe.Field(),
					Message: fe.Error(),
				})
			}
			return false, &ValidationResponse{
				Errors: errors,
			}
		}

		return false, &ValidationResponse{
			Errors: []ValidationError{
				{
					Field:   "",
					Message: err.Error(),
				},
			},
		}
	}

	return true, nil
}

func parseRequestBodyStrict(data any, ctx *fiber.Ctx) error {
	decoder := json.NewDecoder(bytes.NewReader(ctx.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}
