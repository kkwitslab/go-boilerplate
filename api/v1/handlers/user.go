package handlers

import (
	"go-boilerplate/api"
	"go-boilerplate/api/v1/schemas"
	"go-boilerplate/internal/services"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}
func (uh *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var req schemas.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return api.Error{
			Code: http.StatusBadRequest,
			Err:  "failed to parse request body",
		}
	}

	user, err := uh.UserService.CreateUser(req)
	if err != nil {
		return api.Error{
			Code: http.StatusInternalServerError,
			Err:  "failed to create user",
		}
	}

	return api.Response{
		Code: http.StatusOK,
		Msg:  "User Created",
		Data: schemas.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
	}
}
