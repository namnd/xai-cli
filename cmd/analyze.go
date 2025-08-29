/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/namnd/xai-cli/local"
	"github.com/namnd/xai-cli/local/functions"
	"github.com/namnd/xai-cli/xai"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [file_path|directory]",
	Short: "Analyze a file or directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Please provide a file_path\n")
			os.Exit(1)
		}

		_, err := os.Stat(args[0])
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "File not found")
			}
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error openning file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close() // Ensure the file is closed after checking

		err = analyze(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

func analyze(filePath string) error {
	apiKey, err := local.ReadAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	messages := []xai.ChatMessage{
		{
			Role:    "system",
			Content: functions.SystemPrompt,
		},
	}

	fileInfo, _ := os.Stat(filePath)

	var files []string
	if fileInfo.IsDir() {
		files, err = functions.ScanDirectory(filePath)

		if err != nil {
			return fmt.Errorf("failed to scan directory: %v", err)
		}
	} else {
		files = append(files, filePath)
	}

	prompt := "Analyze " + filePath
	messages = append(messages, xai.ChatMessage{
		Role:    "user",
		Content: fmt.Sprintf("Analyze codebase files: %v, then summarize the project. Try to keep it short & concise", files),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	const MAX_ITERATION = 10
	iteration := 0

	threadID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID V7: %v", err)
	}

	for {
		if iteration >= MAX_ITERATION {
			fmt.Println("Max iteration limit reached, exit")
			os.Exit(1)
		}

		iteration++

		chatRequest := xai.ChatRequest{
			Model:      "grok-3-mini",
			Messages:   messages,
			Tools:      functions.Tools,
			ToolChoice: "auto",
		}

		requestBody, err := json.Marshal(chatRequest)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		response, err := xai.MakeAPICall(ctx, apiKey, requestBody)
		if err != nil {
			return fmt.Errorf("failed to make API call: %v", err)
		}

		chatThread, err := local.StoreChat(threadID.String(), prompt, prompt, string(requestBody), string(response))
		if err != nil {
			return fmt.Errorf("failed to store chat history: %v", err)
		}

		var chatResponse xai.ChatResponse
		if err := json.Unmarshal(response, &chatResponse); err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}

		if len(chatResponse.Choices) == 0 {
			fmt.Printf("chat response: %v", chatResponse)
			break
		}

		responseMessage := chatResponse.Choices[0].Message

		if len(responseMessage.ToolCalls) == 0 {
			// completed analyze
			chatThread.ChatRequest = chatRequest
			chatThread.ChatResponse = chatResponse
			chatThread.OriginalPrompt = prompt
			chatThread.Prompt = prompt
			chatThread.ThreadID = threadID.String()
			s, _ := json.Marshal(chatThread)
			fmt.Println(string(s))
			break
		}

		for _, toolCall := range responseMessage.ToolCalls {
			var content string
			switch toolCall.Function.Name {
			case "get_file_content":
				var args functions.FileRequest
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					fmt.Printf("failed to parse tool call arguments: %v", err)
					continue
				}

				result, err := functions.GetFileContent(args.FilePath)
				if err != nil {
					fmt.Printf("failed to execute function: %v", err)
					continue
				}

				resultJSON, _ := json.Marshal(result)
				content = string(resultJSON)
			}

			messages = append(messages, xai.ChatMessage{
				Role:       "tool",
				Content:    content,
				ToolCallID: toolCall.ID,
			})
		}

		// Check context timeout
		select {
		case <-ctx.Done():
			fmt.Printf("Error: Operation timed out: %v\n", ctx.Err())
			os.Exit(1)
		default:
			// Continue loop
		}
	}

	return nil
}
