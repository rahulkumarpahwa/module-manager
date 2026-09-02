package utils

import (
	"os"
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

		if IsProject(dir, entry.Name()) {
			return true
		}
	}

	return false
}
