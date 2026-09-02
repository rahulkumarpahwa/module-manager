package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rahulkumarpahwa/module-manager/services/types"
	"github.com/rahulkumarpahwa/module-manager/services/utils"
)

type Services struct {
	MainFolder *types.MainFolder
}

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

	s.MainFolder.CreationDate, s.MainFolder.ModifiedDate = utils.GetTimes(info)

	if !info.IsDir() {
		return nil, fmt.Errorf("ScanMainFolder: path is not a directory/folder")
	}

	allEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("ScanMainFolder: %w", err)
	}

	entries := paginateEntries(allEntries, offset, limit)

	pesChan := make(chan types.Project)
	wsChan := make(chan types.WorkSpace)

	s.MainFolder.TotalSize = 0
	var wg sync.WaitGroup

	for i, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if entry.Name() == "" {
			continue
		}

		entryPath := filepath.Join(path, entry.Name())

		wg.Add(1)
		go func(entryPath string, entry os.DirEntry, i int) error {
			defer wg.Done()

			if utils.IsProject(entryPath) {

				pe, err := utils.GenerateProjectDetails(entryPath, entry, i)
				if err != nil {
					return fmt.Errorf("ScanMainFolder: %w", err)
				}
				pesChan <- *pe

			} else if utils.IsWorkspace(entryPath) {

				we, err := utils.GenerateWorkspaceDetails(entryPath, entry, i)
				if err != nil {
					return fmt.Errorf("ScanMainFolder: %w", err)
				}
				wsChan <- *we
			}

			return nil
		}(entryPath, entry, i)
	}
	go func() {
		wg.Wait()
		close(pesChan)
		close(wsChan)
	}()

	var pes []types.Project
	var ws []types.WorkSpace

	for pesChan != nil || wsChan != nil {
		select {
		case result, ok := <-pesChan:
			if !ok {
				pesChan = nil
				continue
			}
			pes = append(pes, result)

		case result, ok := <-wsChan:
			if !ok {
				wsChan = nil
				continue
			}
			ws = append(ws, result)
		}
	}

	s.MainFolder.Projects = append(s.MainFolder.Projects, pes...)
	s.MainFolder.ProjectsCount = len(s.MainFolder.Projects)
	s.MainFolder.WorkSpaces = append(s.MainFolder.WorkSpaces, ws...)
	s.MainFolder.WorkSpacesCount = len(s.MainFolder.WorkSpaces)

	for _, project := range s.MainFolder.Projects {
		s.MainFolder.TotalSize += project.Size
	}
	s.MainFolder.IndexProject = makeIndexProject(s.MainFolder.Projects)
	s.MainFolder.IndexWorkspace = makeIndexWorkspace(s.MainFolder.WorkSpaces)

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

func (s *Services) GetEntryById(id string) (interface{}, error) {
	if id == "" {
		return nil, fmt.Errorf("GetEntryById: id cannot be empty.")
	}

	if s.MainFolder == nil {
		return nil, errors.New("GetEntryById: mainfolder cannot be empty.")
	}

	project, ok := s.MainFolder.IndexProject[id]
	if !ok {
		workspace, ok := s.MainFolder.IndexWorkspace[id]
		if !ok {
			return nil, fmt.Errorf("GetEntryById: No Entry (Project or Workspace) with the id: %s", id)
		}
		return &workspace, nil
	}
	return &project, nil
}
