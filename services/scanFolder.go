package services

import (
	"errors"
	"log"
	"module-manager/types"
	"module-manager/utils"
	"os"
	"path/filepath"
)

func MustSetPath(path string) (*types.ScanFolderTypes, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty.")
	}
	return &types.ScanFolderTypes{
		Path:         path,
		TotalSize:    0,
		ProjectCount: 0,
	}, nil
}

func ScanFolder(s *types.ScanFolderTypes) ([]types.ProjectEntry, error) {
	path := filepath.Clean(s.Path)

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

	for _, entry := range projectEntries {
		var pe types.ProjectEntry

		if !entry.IsDir() {
			// when entry is not a Dir we will not consider.
			continue
		}

		var path = filepath.Join(path, entry.Name())

		var size int64
		size, err = utils.GetDirSize(path)
		if err != nil {
			log.Fatal("utils.GetDirSize: " + err.Error())
			return nil, errors.New("utils.GetDirSize: " + err.Error())
		}

		pe = types.ProjectEntry{
			Name:      entry.Name(),
			Path:      path,
			IsDir:     entry.IsDir(),
			TotalSize: size,
		}
		s.TotalSize += size
		s.ProjectCount++
		pes = append(pes, pe)
	}

	return pes, nil
}
