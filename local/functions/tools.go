package functions

import "github.com/namnd/xai-cli/xai"

var SystemPrompt = "You are a code analysis assistant. Use the get_file_content function to retrieve file contents and provide insights about the codebase structure, purpose, and key components. Summarize the code and explain its functionality."

var Tools = []xai.Tool{
	file_content_definition,
}
