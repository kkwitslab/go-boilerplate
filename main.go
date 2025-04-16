package main

import (
	"github.com/gofiber/fiber/v2"
	"go-boilerplate/api/middleware"
	v1 "go-boilerplate/api/v1/routes"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.FiberErrorHandler})

	v1.SetupRoutes(app)

	if err := app.Listen("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}
