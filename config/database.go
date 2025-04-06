package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type DatabaseConfig struct {
	Driver   string
	Filename string
}

func InitDB() error {
	cfg := DatabaseConfig{
		Driver:   "sqlite3",
		Filename: filepath.Join("database", "database.sqlite"),
	}

	if err := os.MkdirAll("database", 0755); err != nil {
		return fmt.Errorf("failed to create the database directory: %v", err)
	}

	var err error
	DB, err = sql.Open(cfg.Driver, cfg.Filename)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}

	log.Println("Database connection established!")

	if err := RunMigrations(); err != nil {
		return fmt.Errorf("failed to run the migrations: %v", err)
	}

	return nil
}

func RunMigrations() error {
	var tableExists int
	err := DB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("error checking the migrations table: %v", err)
	}

	if tableExists == 0 {
		_, err := DB.Exec(`
			CREATE TABLE migrations (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return fmt.Errorf("error creating migrations table: %v", err)
		}
	}

	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("error reading migrations files: %v", err)
	}

	for _, file := range files {
		var migrationRun int
		migrationName := filepath.Base(file)

		err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migrationName).Scan(&migrationRun)
		if err != nil {
			return fmt.Errorf("error checking migrations status: %v", err)
		}

		if migrationRun == 0 {
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("error reading migrations file %s: %v", file, err)
			}

			_, err = DB.Exec(string(content))
			if err != nil {
				return fmt.Errorf("error executing migrations on %s: %v", file, err)
			}

			_, err = DB.Exec("INSERT INTO migrations (name) VALUES (?)", migrationName)
			if err != nil {
				return fmt.Errorf("error recording migrations %s: %v", file, err)
			}

			log.Printf("Applied migration: %s", migrationName)
		}
	}
	return nil
}

func RunSeeders() error {
	files, err := filepath.Glob("database/seeders/*.sql")
	if err != nil {
		return fmt.Errorf("error reading seeder files: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading seeder file %s: %v", file, err)
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error executing seeder %s: %v", file, err)
		}

		log.Printf("Applied seeder: %s", filepath.Base(file))
	}
	return nil
}
