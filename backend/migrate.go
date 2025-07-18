package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := "postgres://sukimise_user:sukimise_password@localhost:5432/sukimise?sslmode=disable"
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) UNIQUE NOT NULL,
			executed_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}

	// Get all migration files
	migrationsDir := "./migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	// Filter and sort up migrations
	var upMigrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			upMigrations = append(upMigrations, file.Name())
		}
	}
	sort.Strings(upMigrations)

	// Execute migrations
	for _, filename := range upMigrations {
		// Check if migration already executed
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = $1", filename).Scan(&count)
		if err != nil {
			log.Fatal("Failed to check migration status:", err)
		}

		if count > 0 {
			fmt.Printf("Migration %s already executed, skipping\n", filename)
			continue
		}

		// Read and execute migration
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			log.Fatal("Failed to read migration file:", err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", filename, err)
		}

		// Record migration as executed
		_, err = db.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
		if err != nil {
			log.Fatal("Failed to record migration:", err)
		}

		fmt.Printf("Migration %s executed successfully\n", filename)
	}

	fmt.Println("All migrations completed successfully")
}
