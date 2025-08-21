package local

import (
	"fmt"
	"os"
	"path/filepath"

	"database/sql"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const (
	TABLE_NAME = "chat_history"
)

func dbPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("failed to get home directory", err)
		os.Exit(1)
	}

	return filepath.Join(homeDir, CONFIG_DIR, CHAT_DB)
}

func InitDB() error {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		thread_id TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		user_message TEXT,
		ai_response TEXT);`, TABLE_NAME)

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func StoreChat(thread_id, user_message, ai_response string) error {
	db, err := sql.Open("sqlite3", dbPath())
	defer db.Close()

	insertChatSQL := fmt.Sprintf(`INSERT INTO %s (thread_id, user_message, ai_response) VALUES (?, ?, ?)`, TABLE_NAME)
	_, err = db.Exec(insertChatSQL, thread_id, user_message, ai_response)
	return err
}
