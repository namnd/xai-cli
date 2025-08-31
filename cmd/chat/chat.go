/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package chat

import (
	"bufio"
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
	"github.com/spf13/viper"
)

// chatCmd represents the chat command
var ChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Interactive chat",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := chat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func chat() error {
	apiKey, err := local.ReadAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	messages := []xai.ChatMessage{
		{
			Role:    "system",
			Content: "You are a highly skilled programming assistant. Provide accurate, concise, and practical solutions for coding tasks. Include code snippets, explanations, and best practices when appropriate. Ask for clarification if the query is ambiguous.",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Starting interactive chat for programming tasks. Type 'exit' to quit.")
	fmt.Print("You: ")

	threadID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID V7: %v", err)
	}

	var i int = 0
	var originalPrompt string

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting chat...")
			break
		}

		if input == "" {
			fmt.Print("You: ")
			continue
		}

		if i == 0 {
			originalPrompt = input
		}

		i++

		messages = append(messages, xai.ChatMessage{
			Role:    "user",
			Content: input,
		})

		chatRequest := xai.ChatRequest{
			Model:    viper.GetString("model"),
			Messages: messages,
		}

		requestBody, err := json.Marshal(chatRequest)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		response, err := xai.MakeAPICall(ctx, apiKey, requestBody)
		if err != nil {
			return fmt.Errorf("failed to make API call: %v", err)
		}

		_, err = local.StoreChat(threadID.String(), originalPrompt, input, string(requestBody), string(response))
		if err != nil {
			fmt.Printf("failed to store chat history: %v", err)
		}

		var chatResponse xai.ChatResponse
		if err := json.Unmarshal(response, &chatResponse); err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}

		if len(chatResponse.Choices) == 0 {
			fmt.Println("No response from API")
			fmt.Print("You:")
			continue
		}

		// Extract assistant response
		assistantMessage := chatResponse.Choices[0].Message.Content
		fmt.Printf("Assistant: %s\n", assistantMessage)

		// Add assistant response to history
		messages = append(messages, xai.ChatMessage{
			Role:    "assistant",
			Content: assistantMessage,
		})

		fmt.Println("====================================")
		fmt.Print("You: ")
	}

	return nil
}
