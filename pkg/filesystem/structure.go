package filesystem

import (
	"os"
	"path/filepath"
)

func CreateProjectStructure(projectPath string, folders []string) error {
	if _, err := os.Stat(projectPath); err == nil {
		if err := os.RemoveAll(projectPath); err != nil {
			return err
		}
	}

	for _, folder := range folders {
		folderPath := filepath.Join(projectPath, folder)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			return err
		}
	}

	return nil
}
