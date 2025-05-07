package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kkwitslab/go-boilerplate/api"
)

func HandleHealthCheck(c *fiber.Ctx) error {

	return api.Response{
		Code: http.StatusOK,
		Msg:  "Service OK",
	}
}
