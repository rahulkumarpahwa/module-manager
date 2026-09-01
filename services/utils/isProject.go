package utils

import (
	"os"
	"path/filepath"
)

func IsProject(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "package.json"))
	return err == nil
}

