package services

import (
	"errors"
	"fmt"
	"log"
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
		return nil, errors.New("path must exist")
	}

	path := filepath.Clean(s.MainFolder.Path)

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("ScanMainFolder: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("ScanMainFolder: path is not a directory/folder")
	}

	Entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("os.ReadDir: " + err.Error())
		return nil, errors.New("os.ReadDir: " + err.Error())
	}

	entries := paginateEntries(Entries, offset, limit)
	pes := make([]types.Project, len(entries))

	s.MainFolder.TotalSize = 0
	var wg sync.WaitGroup

	for i, entry := range entries {
		if !entry.IsDir() {
			// when entry is not a Dir we will not consider.
			continue
		}

		if entry.Name() == "" {
			// No entry without the Name
			continue
		}

		projectPath := filepath.Join(path, entry.Name())

		wg.Add(1)
		go func(projectPath string, entry os.DirEntry, i int) (*types.Project, error) {
			defer wg.Done()
			pe, err := utils.GenerateProjectDetails(projectPath, entry, i)
			if err != nil {
				return nil, errors.New("ScanMainFolder.getProjectDetails: " + err.Error())
			}
			pes[i] = *pe
			return nil, nil
		}(projectPath, entry, i)

	}
	wg.Wait()

	s.MainFolder.ProjectsCount = len(entries)
	s.MainFolder.Projects = append(s.MainFolder.Projects, pes...)
	
	for _, project := range s.MainFolder.Projects {
		s.MainFolder.TotalSize += project.Size
	}
	s.MainFolder.Index = makeIndex(pes)

	return s.MainFolder, nil
}

func paginateEntries(items []os.DirEntry, start int, length int) []os.DirEntry {
	if start >= len(items) {
		return []os.DirEntry{}
	}

	end := start + length
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

func makeIndex(projects []types.Project) map[string]any {
	projectsByID := make(map[string]any, len(projects))

	for _, project := range projects {
		projectsByID[project.ID] = project
	}

	return projectsByID
}

func (s *Services) GetItemInfoById(id string) (types.Project, error) {
	if id == "" {
		return types.Project{}, errors.New("id cannot be empty.")
	}

	if s.MainFolder == nil {
		return types.Project{}, errors.New("scanFolderTypes cannot be empty.")
	}

	project, ok := s.MainFolder.Index[id]
	if !ok {
		return types.Project{}, fmt.Errorf("No project with the id: %s", id)
	}

	return (project), nil
}
