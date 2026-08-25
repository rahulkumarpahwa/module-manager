package services

import (
	"errors"
	"fmt"
	"log"
	"module-manager/types"
	"module-manager/utils"
	"os"
	"path/filepath"
)

func ScanFolder(sft *types.ScanFolderTypes, limit int, offset int) ([]types.ProjectEntry, error) {
	path := filepath.Clean(sft.Path)

	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("folder does not exist.")
		return nil, errors.New("folder does not exist.")
	}

	if !info.IsDir() {
		log.Fatalf("path is not a directory/folder.")
		return nil, errors.New("path is not a directory/folder.")
	}

	projectEntries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("os.ReadDir: " + err.Error())
		return nil, errors.New("os.ReadDir: " + err.Error())
	}

	var pes []types.ProjectEntry

	for i, entry := range take(projectEntries, offset, limit) {

		var pe types.ProjectEntry

		if !entry.IsDir() {
			// when entry is not a Dir we will not consider.
			continue
		}

		var name = entry.Name()
		var path = filepath.Join(path, name)

		var size int64
		size, err = utils.GetDirSize(path)
		if err != nil {
			log.Fatal("utils.GetDirSize: " + err.Error())
			return nil, errors.New("utils.GetDirSize: " + err.Error())
		}

		pe = types.ProjectEntry{
			ID:        fmt.Sprintf("%d%s", i, name),
			Name:      name,
			Path:      path,
			IsDir:     entry.IsDir(),
			TotalSize: size,
		}
		sft.TotalSize += size
		sft.ProjectCount++
		pes = append(pes, pe)
	}

	return pes, nil
}

func take(items []os.DirEntry, start int, length int) []os.DirEntry {
	if start >= len(items) {
		return []os.DirEntry{}
	}

	end := start + length
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}
