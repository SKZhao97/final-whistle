package main

import (
	"flag"
	"log"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/seed"
)

func main() {
	scope := flag.String("scope", "base", "seed scope: base or all")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch *scope {
	case "base":
		if err := seed.SeedBaseData(database.DB); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	case "all":
		if err := seed.SeedAll(database.DB); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	default:
		log.Fatalf("Unknown seed scope: %s", *scope)
	}

	if err := seed.ValidateSeedData(database.DB); err != nil {
		log.Fatalf("Seed validation failed: %v", err)
	}

	log.Printf("%s seed completed successfully", *scope)
}
