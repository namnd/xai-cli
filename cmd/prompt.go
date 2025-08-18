/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runPrompt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPrompt() error {
	apiKey, err := readAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	fmt.Print("Prompt: ")
	reader := bufio.NewReader(os.Stdin)
	userPrompt, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read prompt: %w", err)
	}

	resp, err := makeAPICall(apiKey, userPrompt)
	if err != nil {
		return fmt.Errorf("failed to call xAI API: %w", err)
	}

	fmt.Printf("Response: %s\n", resp)
	
	return nil
}

func readAPIKey() (string, error) {
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

func makeAPICall(apiKey, prompt string) (string, error) {
	client := &http.Client{}
	payload := map[string]any{
		"stream": false,
		"model": "grok-3-mini",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are Grok, a highly intelligent, helpful AI assistant.",
			},
			{
				"role": "user",
				"content": prompt,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON payload: %w", err)
	}

	fmt.Println(bytes.NewBuffer(jsonPayload))
	req, err := http.NewRequest("POST", "https://api.x.ai/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}
