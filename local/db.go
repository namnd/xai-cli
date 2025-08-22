package local

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/namnd/xai-cli/xai"
)

const (
	TABLE_NAME = "chat_history"
)

type ChatThread struct {
	ID           string
	ThreadID     string
	Timestamp    time.Time
	ChatRequest  xai.ChatRequest
	ChatResponse xai.ChatResponse
}

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
		chat_request TEXT,
		chat_response TEXT);`, TABLE_NAME)

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func StoreChat(thread_id, chat_request, chat_response string) error {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	insertChatSQL := fmt.Sprintf(`INSERT INTO %s (thread_id, chat_request, chat_response) VALUES (?, ?, ?)`, TABLE_NAME)
	_, err = db.Exec(insertChatSQL, thread_id, chat_request, chat_response)
	return err
}

func ListRecentChatThread() ([]ChatThread, error) {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	listQuery := fmt.Sprintf(`SELECT thread_id, timestamp, chat_request, chat_response FROM %s GROUP BY thread_id ORDER BY timestamp DESC`, TABLE_NAME)

	rows, err := db.Query(listQuery)

	var chatThreads []ChatThread
	for rows.Next() {
		var chatThread ChatThread
		var chatRequest, chatResponse string
		err := rows.Scan(&chatThread.ID, &chatThread.Timestamp, &chatRequest, &chatResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		err = json.Unmarshal([]byte(chatRequest), &chatThread.ChatRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal chatRequest: %v", err)
		}

		err = json.Unmarshal([]byte(chatResponse), &chatThread.ChatResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal chatResponse: %v", err)
		}

		chatThreads = append(chatThreads, chatThread)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate row: %v", err)
	}

	return chatThreads, nil
}
