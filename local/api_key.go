package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadAPIKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configFile := filepath.Join(homeDir, ".xai", "config")
	apiKey, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	return strings.TrimSpace(string(apiKey)), nil
}
