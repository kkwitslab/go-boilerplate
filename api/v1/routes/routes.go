package routes

import (
	"go-boilerplate/api/middleware"
	"go-boilerplate/api/v1/handlers"
	"go-boilerplate/api/v1/schemas"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	v1 := router.Group("/api/v1")

	setupHealthCheckRoutes(v1)
	setupUserRoutes(v1)
}

func setupHealthCheckRoutes(router fiber.Router) {
	router.Get(
		"/healthz",
		handlers.HandleHealthCheck,
	)
}

func setupUserRoutes(router fiber.Router) {
	user := router.Group("/user")
	userHandler := handlers.NewUserHandler()
	user.Post("/",
		middleware.ValidatorMiddleware[schemas.CreateUserRequest](),
		userHandler.HandleCreateUser,
	)
}
