package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kkwitslab/go-boilerplate/api"
)

// FiberErrorHandler Fiber error handler for app.Config.ErrorHandler
func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	var apiRes api.Response
	var apiErr api.Error

	switch {
	case errors.As(err, &apiRes):
		return ctx.Status(apiRes.Code).JSON(apiRes)

	case errors.As(err, &apiErr):
		return ctx.Status(apiErr.Code).JSON(apiErr)

	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code: fiber.StatusInternalServerError,
			Err:  "Internal Server Error",
		})
	}
}
