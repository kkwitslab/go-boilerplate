package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-boilerplate/api"
	"go-boilerplate/api/v1/schemas"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidatorMiddleware[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data T

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		if err := validate.Struct(data); err != nil {
			var fieldErrors []schemas.FieldError
			for _, fe := range err.(validator.ValidationErrors) {
				fieldErrors = append(fieldErrors, schemas.FieldError{
					Code:    "VALIDATION_ERROR",
					Field:   fe.Field(),
					Message: generateErrorMessage(fe),
					Data:    fe.Tag(),
				})
			}

			return api.Error{
				Code: http.StatusBadRequest,
				Err:  "Invalid request payload",
				Data: fieldErrors,
			}
		}

		c.Locals("validatedBody", data)
		return c.Next()
	}
}

func generateErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "email":
		return fe.Field() + " must be a valid email"
	default:
		return fe.Field() + " is invalid"
	}
}
