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

	"github.com/namnd/xai-cli/local"
	"github.com/namnd/xai-cli/xai"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze codebase in the current working directory",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := analyze()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

func analyze() error {
	apiKey, err := local.ReadAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	tools := []xai.Tool{
		{
			Type: "function",
			Function: xai.FunctionDetails{
				Name:        "get_file_content",
				Description: "Retrieve the content of a specific file in the codebase",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file_path": map[string]any{
							"type":        "string",
							"description": "Path to the file in the codebase",
						},
					},
					"required": []string{"file_path"},
				},
			},
		},
	}

	messages := []xai.FunctionCallMessage{
		{
			Role:    "system",
			Content: "You are a code analysis assistant. Use the get_file_content function to retrieve file contents and provide insights about the codebase structure, purpose, and key components. Summarize the code and explain its functionality.",
		},
		{
			Role:    "user",
			Content: "Analyze the main.go file in my codebase.",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	const MAX_ITERATION = 10
	iteration := 0

	for {
		if iteration >= MAX_ITERATION {
			fmt.Println("max iteration limit reached, exit")
			os.Exit(1)
		}

		iteration++

		chatRequest := xai.FunctionCallRequest{
			Model:      "grok-3-mini",
			Messages:   messages,
			Tools:      tools,
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

		var chatResponse xai.FunctionCallResponse
		if err := json.Unmarshal(response, &chatResponse); err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}

		if len(chatResponse.Choices) == 0 {
			fmt.Printf("chat response: %v", chatResponse)
			break
		}

		responseMessage := chatResponse.Choices[0].Message

		if len(responseMessage.ToolCalls) == 0 {
			fmt.Printf("assistant response: %s", responseMessage.Content)
			break
		}

		for _, toolCall := range responseMessage.ToolCalls {
			var content string
			switch toolCall.Function.Name {
			case "get_file_content":
				var args local.FileRequest
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					fmt.Printf("failed to parse tool call arguments: %v", err)
					continue
				}

				result, err := local.GetFileContent(args.FilePath)
				if err != nil {
					fmt.Println("failed to execute function: %w", err)
					continue
				}

				resultJSON, _ := json.Marshal(result)
				content = string(resultJSON)


			}
			messages = append(messages, xai.FunctionCallMessage{
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
