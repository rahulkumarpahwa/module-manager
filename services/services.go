package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	// "sync"

	"github.com/rahulkumarpahwa/module-manager/services/types"
	"github.com/rahulkumarpahwa/module-manager/services/utils"
)

type Services struct {
	MainFolder *types.MainFolder
}

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 50
)

func NewServices(mainFolder *types.MainFolder) *Services {
	return &Services{
		MainFolder: mainFolder,
	}
}

func (s *Services) ScanMainFolder(limit int, offset int) (*types.MainFolder, error) {
	if s.MainFolder.Path == "" {
		return nil, fmt.Errorf("ScanMainFolder: s.MainFolder.Path must exist")
	}

	path := filepath.Clean(s.MainFolder.Path)

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("ScanMainFolder: %w", err)
	}

	s.MainFolder.TotalSize = info.Size()
	s.MainFolder.CreationDate, s.MainFolder.ModifiedDate = utils.GetTimes(info)

	if !info.IsDir() {
		return nil, fmt.Errorf("ScanMainFolder: path is not a directory/folder")

	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ScanMainFolder: %w", err)
	}
	defer f.Close()

	if offset < defaultOffset {
		offset = defaultOffset
	}

	if limit <= 0 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	if offset > 0 {
		_, err = f.Readdir(offset)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("skip entries: %w", err)
		}
	}

	entries, err := f.Readdir(limit)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read entries: %w", err)
	}

	s.MainFolder.FirstRunEntryCount = len(entries)

	for _, entry := range entries {
		if !entry.IsDir() {
			return s.MainFolder, nil
		}

		var firstRunEntry types.FirstRunEntry
		firstRunEntry.ID = utils.GenerateID()
		firstRunEntry.Name = entry.Name()
		firstRunEntry.Path = filepath.Join(path, entry.Name())
		firstRunEntry.ModifiedDate = entry.ModTime()

		detectInfo := utils.DetectProject(firstRunEntry.Path)
		firstRunEntry.Language = string(detectInfo.Language)
		firstRunEntry.IsProject = detectInfo.IsProject
		firstRunEntry.IsWorkspace = !firstRunEntry.IsProject

		s.MainFolder.FirstRunEntries = append(s.MainFolder.FirstRunEntries, firstRunEntry)
	}
	
	return s.MainFolder, nil
}

func paginateEntries(items []os.DirEntry, start int, length int) []os.DirEntry {
	if start < 0 {
		return []os.DirEntry{}
	}

	if start >= len(items) {
		return []os.DirEntry{}
	}

	if length < 0 {
		return []os.DirEntry{}
	}

	end := start + length
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

func makeIndexProject(entries []types.Project) map[string]types.Project {
	entriesByID := make(map[string]types.Project, len(entries))

	for _, entry := range entries {
		entriesByID[entry.ID] = entry
	}

	return entriesByID
}

func makeIndexWorkspace(entries []types.WorkSpace) map[string]types.WorkSpace {
	entriesByID := make(map[string]types.WorkSpace, len(entries))

	for _, entry := range entries {
		entriesByID[entry.ID] = entry
	}

	return entriesByID
}

// func (s *Services) GetEntryById(id string) (interface{}, error) {
// 	if id == "" {
// 		return nil, fmt.Errorf("GetEntryById: id cannot be empty.")
// 	}

// 	if s.MainFolder == nil {
// 		return nil, errors.New("GetEntryById: mainfolder cannot be empty.")
// 	}

// 	project, ok := s.MainFolder.IndexProject[id]
// 	if !ok {
// 		workspace, ok := s.MainFolder.IndexWorkspace[id]
// 		if !ok {
// 			return nil, fmt.Errorf("GetEntryById: No Entry (Project or Workspace) with the id: %s", id)
// 		}
// 		return &workspace, nil
// 	}
// 	return &project, nil
// }
