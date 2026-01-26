package main

import (
	"log"

	"example-app/internal/api/v1"
	"example-app/internal/models"
	"github.com/project/go-framework/pkg/app"
)

func main() {
	a, err := app.New("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Migrate database
	if err := a.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Register routes
	v1.RegisterRoutes(a.Server.Engine, a.DB, a.Logger)

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
