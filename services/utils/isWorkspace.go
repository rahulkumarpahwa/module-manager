package utils

import (
	"os"
	"path/filepath"
)

func IsWorkspace(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		child := filepath.Join(dir, entry.Name())

		if IsProject(child) {
			return true
		}
	}

	return false
}
