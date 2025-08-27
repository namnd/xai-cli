package functions

import (
	"io/fs"
	"path/filepath"
	"slices"
)

var importantExts = []string{".go", ".js", ".py", "README.md"}

func ScanDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() { // Only files, not directories
			if slices.Contains(importantExts, filepath.Ext(path)) {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}
