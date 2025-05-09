package di

import (
	"fmt"
	"log"

	"github.com/kkwitslab/go-boilerplate/internal/config"
	"github.com/kkwitslab/go-boilerplate/internal/repositories"
	"github.com/kkwitslab/go-boilerplate/internal/services"

	"go.uber.org/dig"

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

func init() {
	if err := initializeDeps(); err != nil {
		log.Fatalf("error while initializing app dependencies: %v\n", err)
	}

	return
}
