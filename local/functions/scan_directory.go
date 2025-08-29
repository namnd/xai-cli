package functions

import (
	"io/fs"
	"path/filepath"
)

func ScanDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() { // Only files, not directories
			if _, found := SupportedLanguages[filepath.Ext(path)]; found {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}
