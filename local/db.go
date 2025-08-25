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

func (c *ChatThread) Display() {
	s, _ := json.Marshal(c)
	fmt.Println(string(s))
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

func StoreChat(thread_id, chat_request, chat_response string) (*ChatThread, error) {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	insertChatSQL := fmt.Sprintf(`INSERT INTO %s (thread_id, chat_request, chat_response) VALUES (?, ?, ?)`, TABLE_NAME)
	_, err = db.Exec(insertChatSQL, thread_id, chat_request, chat_response)

	return &ChatThread{
		ThreadID: thread_id,
	}, nil

}

func GetThreadByID(threadID string) (*ChatThread, error) {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	getThreadQuery := fmt.Sprintf(`SELECT timestamp, chat_request, chat_response FROM %s WHERE thread_id = ? ORDER BY timestamp DESC LIMIT 1`, TABLE_NAME)

	row := db.QueryRow(getThreadQuery, threadID)
	var chatThread ChatThread
	var chatRequest, chatResponse string

	err = row.Scan(&chatThread.Timestamp, &chatRequest, &chatResponse)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("threadID %s not found", threadID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan chat history: %v", err)
	}

	err = json.Unmarshal([]byte(chatRequest), &chatThread.ChatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chatRequest: %v", err)
	}

	err = json.Unmarshal([]byte(chatResponse), &chatThread.ChatResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chatResponse: %v", err)
	}

	chatThread.ThreadID = threadID

	return &chatThread, nil
}

func GetChatMinusT(t int) (*ChatThread, error) {
	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	getThreadQuery := fmt.Sprintf(`SELECT thread_id FROM %s GROUP BY thread_id ORDER BY timestamp DESC LIMIT 1 OFFSET ?`, TABLE_NAME)
	row := db.QueryRow(getThreadQuery, t)

	var threadID string
	err = row.Scan(&threadID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("chat thread T-%d not found", t)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan chat thread: %v", err)
	}

	getQuery := fmt.Sprintf(`SELECT timestamp, chat_request, chat_response FROM %s WHERE thread_id = ? ORDER BY timestamp DESC LIMIT 1`, TABLE_NAME)

	row = db.QueryRow(getQuery, threadID)

	var chatThread ChatThread
	var chatRequest, chatResponse string

	err = row.Scan(&chatThread.Timestamp, &chatRequest, &chatResponse)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("chat T-%d not found", t)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan chat history: %v", err)
	}

	err = json.Unmarshal([]byte(chatRequest), &chatThread.ChatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chatRequest: %v", err)
	}

	err = json.Unmarshal([]byte(chatResponse), &chatThread.ChatResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chatResponse: %v", err)
	}

	chatThread.ThreadID = threadID

	return &chatThread, nil
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
