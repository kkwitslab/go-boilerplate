package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-boilerplate/api"
	"net/http"
)

func HandleHealthCheck(c *fiber.Ctx) error {

	return api.Response{
		Code: http.StatusOK,
		Msg:  "Service OK",
	}
}
