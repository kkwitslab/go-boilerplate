package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	v1 "github.com/kkwitslab/go-boilerplate/api/rest/v1"
)

func HandleHealthCheck(c *fiber.Ctx) error {

	return v1.Response{
		Code: http.StatusOK,
		Msg:  "Service OK",
	}
}
