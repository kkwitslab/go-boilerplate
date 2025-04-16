package di

import (
	"fmt"
	"go-boilerplate/api/middleware"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/models"
	"go-boilerplate/internal/repositories"
	"go-boilerplate/internal/services"
	"log"

	"go.uber.org/dig"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Container = dig.New()
)

func initDB() (*gorm.DB, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.AppConfig.DatabaseUser,
		config.AppConfig.DatabasePassword,
		config.AppConfig.DatabaseHost,
		config.AppConfig.DatabasePort,
		config.AppConfig.DatabaseName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (*fiber.App, error) {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.FiberErrorHandler})

	// initialize application dependencies
	if err := initializeDeps(); err != nil {
		return nil, fmt.Errorf("failed to initialize app dependencies: %w", err)
	}

	if err := Container.Invoke(func(db *gorm.DB) {
		// Run database migrations
		err := models.RunMigrations(db)
		if err != nil {
			fmt.Printf("failed to run migrations: %v\n", err)
			return
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to invoke db dependency: %v", err)
	}

	return app, nil
}

func initializeDeps() error {
	db, err := initDB()
	if err != nil {
		log.Fatalf("error initializing the database: %v\n", err)
	}

	// Provide DB interface
	if err := Container.Provide(func() *gorm.DB {
		return db
	}); err != nil {
		return err
	}

	userRepository := repositories.NewPostgresUserRepository(db)
	if err := Container.Provide(func() repositories.UserRepository {
		return userRepository
	}); err != nil {
		return err
	}

	userService := services.NewUserService(userRepository)
	if err := Container.Provide(func() *services.UserService {
		return userService
	}); err != nil {
		return err
	}

	return nil
}
