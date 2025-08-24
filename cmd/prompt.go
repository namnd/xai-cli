/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/namnd/xai-cli/local"
	"github.com/namnd/xai-cli/xai"
	"github.com/spf13/cobra"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt [something]",
	Short: "Enter a prompt to get response from xAI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		userPrompt := strings.Join(args, " ")
		userPrompt = strings.TrimSpace(userPrompt)
		if userPrompt == "" {
			fmt.Fprintln(os.Stderr, "Please provide a valid prompt")
			os.Exit(1)
		}

		err := runPrompt(userPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}

func runPrompt(userPrompt string) error {
	apiKey, err := local.ReadAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	messages := []xai.ChatMessage{
		{
			Role:    "system",
			Content: "You are a highly skilled programming assistant. Provide accurate, concise, and practical solutions for coding tasks. Include code snippets, explanations, and best practices when appropriate. Ask for clarification if the query is ambiguous.",
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	chatRequest := xai.ChatRequest{
		Model:    "grok-3-mini",
		Messages: messages,
	}

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	response, err := xai.MakeAPICall(ctx, apiKey, requestBody)

	var chatResponse xai.ChatResponse
	if err := json.Unmarshal(response, &chatResponse); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return fmt.Errorf("No response from API")
	}

	fmt.Println()
	fmt.Println(chatResponse.Choices[0].Message.Content)

	return nil
}
