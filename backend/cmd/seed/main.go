package main

import (
	"log"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/seed"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := seed.SeedBaseData(database.DB); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	if err := seed.ValidateSeedData(database.DB); err != nil {
		log.Fatalf("Seed validation failed: %v", err)
	}

	log.Println("Base seed completed successfully")
}
