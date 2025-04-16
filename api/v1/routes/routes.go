package routes

import (
	"go-boilerplate/api/middleware"
	"go-boilerplate/api/v1/handlers"
	"go-boilerplate/api/v1/schemas"
	"go-boilerplate/internal/di"
	"go-boilerplate/internal/services"

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
