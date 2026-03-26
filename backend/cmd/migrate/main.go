package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
)

func main() {
	dir := flag.String("dir", "migrations", "directory containing SQL migration files")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	files, err := filepath.Glob(filepath.Join(*dir, "*.sql"))
	if err != nil {
		log.Fatalf("Failed to scan migrations: %v", err)
	}
	if len(files) == 0 {
		log.Fatalf("No SQL migrations found in %s", *dir)
	}

	sort.Strings(files)

	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration %s: %v", file, err)
		}

		log.Printf("Applying migration %s", filepath.Base(file))
		if err := database.DB.Exec(string(sqlBytes)).Error; err != nil {
			log.Fatalf("Failed to apply migration %s: %v", filepath.Base(file), err)
		}
	}

	log.Printf("Applied %d migration files successfully", len(files))
}
