package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"test-project/utils"
)

type DocEntry struct {
	Name  string
	Path  string
	Size  int64
	IsDir bool
}

func main() {
	// We will turn this into a worker and then use again and again.

	var totalSize int64
	path := flag.String("path", "", "folder path")

	flag.Parse()

	if *path == "" {
		log.Fatal("Path can't be empty")
	}

	entries, err := os.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	var details []DocEntry

	for _, entry := range entries {
		var fullPath = filepath.Join(*path, entry.Name())

		var size int64

		if entry.IsDir() {
			size, err = utils.GetDirSize(fullPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				fmt.Println(err)
				continue
			}

			size = info.Size()
		}

		docEntry := DocEntry{
			Name:  entry.Name(),
			Path:  fullPath,
			Size:  size,
			IsDir: entry.IsDir(),
		}

		totalSize += docEntry.Size
		details = append(details, docEntry)
	}

	// Header
	fmt.Printf("%-30s %-50s %-12s %-10s\n",
		"NAME", "PATH", "SIZE", "TYPE")
	fmt.Println("-----------------------------------------------------------------------------------------------")

	// Rows
	for _, entry := range details {
		entryType := "FILE"
		if entry.IsDir {
			entryType = "DIR"
		}

		fmt.Printf("%-30s %-50s %-12s %-10s\n",
			entry.Name,
			entry.Path,
			utils.FormatSize(entry.Size),
			entryType,
		)
	}

	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Printf("Total Size: %s\n", utils.FormatSize(totalSize))
}
