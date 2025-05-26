package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kkwitslab/go-boilerplate/api/rest/v1/middleware"
	"github.com/kkwitslab/go-boilerplate/internal/config"
	"github.com/kkwitslab/go-boilerplate/internal/di"
	"github.com/kkwitslab/go-boilerplate/internal/models"
	"github.com/kkwitslab/go-boilerplate/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

// InitializeApp initializes the application with all dependencies
func InitializeServer() (*fiber.App, error) {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.FiberErrorHandler})

	// initialize application dependencies
	if err := di.Container.Invoke(func(db *gorm.DB) {
		// Run database migrations
		err := models.RunMigrations(db)
		if err != nil {
			fmt.Printf("failed to run migrations: %v\n", err)
			return
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to invoke db dependency: %v", err)
	}

	app.All("/grpc/*", adaptor.HTTPHandler(createGrpcGatewayMux(context.Background())))
	return app, nil
}

func createGrpcGatewayMux(ctx context.Context) http.Handler {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, mux, config.AppConfig.GetGRPCListenAddress(), opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway handler: %v", err)
	}
	return mux
}
