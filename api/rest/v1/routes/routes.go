package routes

import (
	"github.com/kkwitslab/go-boilerplate/api/rest/v1/handlers"
	"github.com/kkwitslab/go-boilerplate/api/rest/v1/middleware"
	"github.com/kkwitslab/go-boilerplate/api/rest/v1/schemas"
	"github.com/kkwitslab/go-boilerplate/internal/di"
	"github.com/kkwitslab/go-boilerplate/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) error {
	v1 := router.Group("/api/v1")

	if err := setupHealthCheckRoutes(v1); err != nil {
		return err
	}

	if err := setupUserRoutes(v1); err != nil {
		return err
	}

	return nil
}

func setupHealthCheckRoutes(router fiber.Router) error {
	router.Get(
		"/healthz",
		handlers.HandleHealthCheck,
	)
	return nil
}

func setupUserRoutes(router fiber.Router) error {
	userPrefix := router.Group("/user")

	var userService *services.UserService
	if err := di.Container.Invoke(func(u *services.UserService) {
		userService = u
	}); err != nil {
		return err
	}

	userHandlers := handlers.NewUserHandler(userService)

	userPrefix.Post("/",
		middleware.ValidatorMiddleware[schemas.CreateUserRequest](),
		userHandlers.HandleCreateUser,
	)

	return nil
}
