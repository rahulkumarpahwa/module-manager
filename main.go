package main

import (
	"flag"
	"log"
	"module-manager/services"
	"module-manager/utils"
)

func main() {
	// We will turn this into a worker and then use again and again.
	path := flag.String("path", "", "folder path")

	flag.Parse()

	if *path == "" {
		log.Fatal("Path can't be empty")
	}

	scanFolder, err := services.MustSetPath(*path)
	if err != nil {
		log.Fatalf("services.MustSetPath: %v", err.Error())
	}

	details, err := services.ScanFolder(scanFolder)
	if err != nil {
		log.Fatalf("Main.services.ScanFolder: %v", err.Error())
	}

	utils.ShowTables(details, scanFolder.TotalSize, scanFolder.ProjectCount)
}
