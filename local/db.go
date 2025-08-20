package local

import (
	"fmt"
	"os"
	"path/filepath"

	"database/sql"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	dbPath := filepath.Join(homeDir, ".xai", "chat_history.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS chat_history(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	user_message TEXT,
	ai_response TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func SaveChat(userMessage, aiResponse string) error {
	insertSQL := `INSERT INTO chat_history (user_message, ai_response) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, userMsg, aiResponse)
	return err
}
