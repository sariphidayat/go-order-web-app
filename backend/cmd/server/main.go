package main

import (
	"log"

	"customer-order-app/backend/internal/config"
	"customer-order-app/backend/internal/database"
	"customer-order-app/backend/internal/routes"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := database.SeedAdmin(db, cfg); err != nil {
		log.Fatalf("admin seed failed: %v", err)
	}

	router := routes.Setup(db, cfg)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
