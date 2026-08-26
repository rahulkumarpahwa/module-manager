package utils

import (
	"io/fs"
	"path/filepath"
)

// In go, we can't have the total folder size directly, we need to add the size of the each file of directory to get total.
type DirSize struct {
	Bytes int64   `json:"bytes"`
	KB    float64 `json:"kb"`
	MB    float64 `json:"mb"`
	GB    float64 `json:"gb"`
}

func GetDirSize(root string) (DirSize, error) {
	var bytes int64

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		bytes += info.Size()
		return nil
	})

	if err != nil {
		return DirSize{}, err
	}

	return DirSize{
		Bytes: bytes,
		KB:    float64(bytes) / 1024,
		MB:    float64(bytes) / (1024 * 1024),
		GB:    float64(bytes) / (1024 * 1024 * 1024),
	}, nil
}
