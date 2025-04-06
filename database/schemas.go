package database

import (
    "database/sql"
)

type User struct {
    ID      int
    Name    string
    Email   string
}

func CreateUseerTable(db *sql.DB) error {
    _, err := db.Exec(`
          CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP    
        )
        `)
    return err
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
