package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

var DB sql.DB

type DatabaseConfig struct {
    Driver  string
    Filename  string
}

func InitDB()  {
    cfg := DatabaseConfig{
        Driver: "sqlite3",
        Filename: filepath.Join("database", "database.sqlite"),
    }

    if err := os.MkdirAll("database", 0755); err != nil {
        log.Fatalf("Failed to create the database directory: %v", err)
    }

    var err error
    DB, err := sql.Open(cfg.Driver, cfg.Filename)
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatalf("Failed to ping the database: %v", err)
    }

    fmt.Printf("Database connection established!")

    if err := RunMigrations(); err != nil {
        log.Fatalf("Failed to run the migrations: %v", err) 
    }
}

func RunMigrations() error{

    var tableExits int
    err := DB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&tableExits)
    if err != nil {
        return fmt.Errorf("Error checking the migrations table: %v", err)
    }

    if tableExits == 0 {
        _, err := DB.Exec(`
               CREATE TABLE migrations (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
            `)
        if err != nil {
            return fmt.Errorf("Error creating migrations table: %v", err)
        }
    }

    files, err := filepath.Glob("database/migrations/*.sql")
    if err != nil {
        return fmt.Errorf("Error reading migrations files: %v", err)
    }

    for _, file := range files{
        var migrationRun int
        migrationName := filepath.Base(file)

        err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migrationName).Scan(migrationRun)
        if err != nil {
            return fmt.Errorf("Error checking migrations status: %v", err)
        }

        if migrationRun == 0 {
            
            content, err := os.ReadFile(file)

            if err != nil {
                return fmt.Errorf("Error reading migrations files %s: %v", file, err)
            }

            _, err = DB.Exec(string(content))
            if err != nil {
                return fmt.Errorf("Error executing migrations on %s: %v", file, err)
            }

            _, err = DB.Exec("INSERT INTO migrations (name) VALUES (?)", migrationName)

            if err != nil {
                return fmt.Errorf("Error recording migrations %s: %v", file, err)
            }

            log.Printf("Applied migrations: %s", migrationName)
        }
    }
    return nil
}

func RunSeeders() error {
	// Read seeder files from database/seeders
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
