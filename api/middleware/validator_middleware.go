package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-boilerplate/api"
	"go-boilerplate/utils"
	"net/http"
)

func ValidatorMiddleware[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data T

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		if err := utils.Validate(&data); err != nil {
			return api.Error{
				Code: http.StatusBadRequest,
				Err:  "Invalid Request Payload",
				Data: err,
			}
		}

		return c.Next()
	}
}
