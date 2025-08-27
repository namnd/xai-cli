package cmd

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"os"

	"github.com/namnd/xai-cli/local"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// setupCmd represents the auth command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Authenticate by entering xAI API key",
	Long: `API key will be stored in ~/.xai/config.

TODO: improve security`,
	Run: func(cmd *cobra.Command, args []string) {
		err := configureAPIKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func configureAPIKey() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".xai")
	configFile := filepath.Join(configDir, "config")

	fmt.Print("Enter your xAI API Key: ")
	apiKey, err := readSecret()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	err = os.WriteFile(configFile, []byte(apiKey), 0600)
	if err != nil {
		return fmt.Errorf("failed to write API key to file: %w", err)
	}

	fmt.Println("xAI API key successfully stored in", configFile)

	if err := local.InitDB(); err != nil {
		fmt.Printf("failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Initialize db successfully")

	return nil
}

func readSecret() (string, error) {
	if !term.IsTerminal(int(syscall.Stdin)) {
		// fallback to regular input
		reader := bufio.NewReader(os.Stdin)
		return reader.ReadString('\n')
	}

	byteInput, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(byteInput), nil
}
