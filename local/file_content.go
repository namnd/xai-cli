package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileRequest struct {
	FilePath string `json:"file_path" description:"Path to the file in the codebase"`
}

type FileContent struct {
	FilePath string `json:"file_path" description:"Path to the file"`
	Content  string `json:"content" description:"Content of the file"`
	Language string `json:"language" description:"Programming language of the file"`
}

func GetFileContent(filePath string) (FileContent, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return FileContent{
				FilePath: filePath,
				Content:  "File not found",
				Language: "unknown",
			}, nil
		}

		return FileContent{
			FilePath: filePath,
			Content:  fmt.Sprintf("Error reading file: %v", err),
			Language: "unknown",
		}, nil
	}

	// Determine language based on file extension
	extension := strings.ToLower(filepath.Ext(filePath))
	languageMap := map[string]string{
		".go":   "go",
		".py":   "python",
		".js":   "javascript",
		".ts":   "typescript",
		".java": "java",
		".cpp":  "cpp",
		".c":    "c",
		".cs":   "csharp",
		".rb":   "ruby",
	}
	language, ok := languageMap[extension]
	if !ok {
		language = "unknown"
	}

	return FileContent{
		FilePath: filePath,
		Content:  string(content),
		Language: language,
	}, nil
}
