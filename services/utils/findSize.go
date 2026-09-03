package utils

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// returns the DirectorySize in bytes (int64)
func GetDirSizeInBytes(ctx context.Context, root string) (int64, error) {
	var bytes int64

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err != nil {
			return err
		}

		if entry.IsDir() && entry.Name() == "node_modules" {
			return fs.SkipDir
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		bytes += info.Size()
		return nil
	})

	if err != nil {
		return 0, err
	}

	return bytes, nil
}

// formats the size to string. takes in size in bytes to format.
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// get the node_modules size in bytes (use workers)
func GetNodeModulesSizeInBytes(ctx context.Context, nodeModulesPath string, workers int) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	files := make(chan string)
	var wg sync.WaitGroup
	// Each worker calculates its own local total.
	var mu sync.Mutex
	var total int64 // Start workers.
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			var localTotal int64
			for path := range files {
				info, err := os.Stat(path)
				if err != nil {
					continue
				}
				localTotal += info.Size()
			}
			mu.Lock()
			total += localTotal
			mu.Unlock()
		}()
	} // Walk the directory and send files to workers.

	err := filepath.WalkDir(nodeModulesPath, func(path string, d fs.DirEntry, err error) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err != nil {
			return err
		}
		if !d.IsDir() {
			files <- path
		}
		return nil
	}) // No more files will be sent.
	close(files) // Wait for all workers.
	wg.Wait()
	return total, err
}

/**

                 WalkDir
                    │
             finds a file
                    │
                    ▼
              files channel
             /      |      \
            ▼       ▼       ▼
        Worker 1 Worker 2 Worker 3 ...
            │       │       │
            ▼       ▼       ▼
        localTotal localTotal localTotal
             \       |       /
              \      |      /
               ▼     ▼     ▼
                  total

				  **/
