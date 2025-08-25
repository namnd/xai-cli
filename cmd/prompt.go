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

	"github.com/google/uuid"
	"github.com/namnd/xai-cli/local"
	"github.com/namnd/xai-cli/xai"
	"github.com/spf13/cobra"
)

var threadID string

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt [something]",
	Short: "Enter a prompt to get response from xAI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		userPrompt := strings.Join(args, " ")
		userPrompt = strings.TrimSpace(userPrompt)
		userPrompt = strings.Trim(userPrompt, "\n")
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

	promptCmd.Flags().StringVarP(&threadID, "thread-id", "t", "", "Continue prompt of the given thread ID. If not provided, prompt will start a new thread")
}

func runPrompt(userPrompt string) error {
	apiKey, err := local.ReadAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var messages []xai.ChatMessage
	if threadID != "" {
		thread, err := local.GetThreadByID(threadID)
		if err != nil {
			return fmt.Errorf("failed to get threadByID: %v", err)
		}
		messages = thread.ChatRequest.Messages
	} else {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("failed to generate UUID V7: %v", err)
		}

		threadID = id.String()
		messages = []xai.ChatMessage{
			{
				Role:    "system",
				Content: "You are a highly skilled programming assistant. Provide accurate, concise, and practical solutions for coding tasks. Include code snippets, explanations, and best practices when appropriate. Ask for clarification if the query is ambiguous.",
			},
		}
	}

	messages = append(messages, xai.ChatMessage{
		Role:    "user",
		Content: userPrompt,
	})

	chatRequest := xai.ChatRequest{
		Model:    "grok-3-mini",
		Messages: messages,
	}

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	response, err := xai.MakeAPICall(ctx, apiKey, requestBody)

	chatThread, err := local.StoreChat(threadID, string(requestBody), string(response))
	if err != nil {
		fmt.Printf("failed to store prompt: %v", err)
	}

	var chatResponse xai.ChatResponse
	if err := json.Unmarshal(response, &chatResponse); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	chatThread.ChatRequest = chatRequest
	chatThread.ChatResponse = chatResponse

	s, _ := json.Marshal(chatThread)
	fmt.Println(string(s))

	return nil
}
