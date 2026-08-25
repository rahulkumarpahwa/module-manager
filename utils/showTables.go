package utils

import (
	"fmt"
	"module-manager/types"
)

func ShowTables(details []types.ProjectEntry, TotalSize int64, ProjectCount int) {
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
			FormatSize(entry.TotalSize),
			entryType,
		)
	}

	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Printf("Total Size: %s\n", FormatSize(TotalSize))
	fmt.Printf("Total Projects Count: %d\n", ProjectCount)
}
