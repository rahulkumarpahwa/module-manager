package utils

import (
	"os"
	"path/filepath"
)

// In go, we can't have the total folder size directly, we need to add the size of the each file of directory to get total.
func GetDirSize(path string) (int64, error) {
	var total int64

	err := filepath.WalkDir(path, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Don't add directories themselves
		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		total += info.Size()
		return nil
	})

	return total, err
}
