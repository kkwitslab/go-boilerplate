package main

import (
	"fmt"

	v1 "github.com/kkwitslab/go-boilerplate/api/v1/routes"
	"github.com/kkwitslab/go-boilerplate/internal/config"
	"github.com/kkwitslab/go-boilerplate/internal/di"

	"log"
)

func main() {
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := v1.SetupRoutes(app); err != nil {
		log.Fatalf("error setting application routes: %v\n", err)
	}

	addr := fmt.Sprintf("%s:%d",
		config.AppConfig.HTTPListenAddress,
		config.AppConfig.HTTPListenPort,
	)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
