package services

import (
	"errors"
	"fmt"
	"github.com/rahulkumarpahwa/module-manager/types"
	"github.com/rahulkumarpahwa/module-manager/utils"
	"log"
	"os"
	"path/filepath"
)

type Services struct {
	MainFolderInfo *types.MainFolderInfo
	Project        *types.Project
}

func NewServices(mainFolderInfo *types.MainFolderInfo) *Services {
	return &Services{
		MainFolderInfo: mainFolderInfo,
	}
}

func (s *Services) ScanMainFolder(limit int, offset int) ([]types.Project, error) {
	if s.MainFolderInfo.Path == "" {
		return nil, errors.New("path must exist.")
	}

	path := filepath.Clean(s.MainFolderInfo.Path)

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

	var pes []types.Project

	entries := take(projectEntries, offset, limit)

	for i, entry := range entries {

		var pe types.Project

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

		pe = types.Project{
			ID:        fmt.Sprintf("%d%s", i+1, name),
			Name:      name,
			Path:      path,
			IsDir:     entry.IsDir(),
			TotalSize: size,
		}
		s.MainFolderInfo.TotalSize += size
		pes = append(pes, pe)
	}

	s.MainFolderInfo.ProjectCount = len(entries)
	s.MainFolderInfo.Projects = pes
	s.MainFolderInfo.Index = makeIndex(pes)

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

func makeIndex(projects []types.Project) map[string]types.Project {
	projectsByID := make(map[string]types.Project, len(projects))

	for _, project := range projects {
		projectsByID[project.ID] = project
	}

	return projectsByID
}

/*
entries, err := os.ReadDir(s.Project.Path)
	if err != nil {
		return nil, errors.New("GetProjectInfo: os.ReadDir: " + err.Error())
	}

	project := &types.Project{}

	for _, entry := range entries {
		var name = entry.Name()

		if name == "package.json" {
			project.HasPackageJson = true
		}

		if name == "node_modules" {
			project.HasNodeModules = true
		}

		if name == ".git" {
			project.HasGit = true
		}

		project.ListFiles = append(project.ListFiles, name)
	}

	if project.HasNodeModules {
		path := filepath.Join(project.Path, "node_modules")
		size, err := utils.GetDirSize(path)
		if err != nil {
			return nil, errors.New("GetProjectInfo: utils.GetDirSize: " + err.Error())
		}
		project.NodeModulesSize = size
	}

	if project.HasGit {
		project.LatestCommit, project.LatestCommitDate = utils.GetLatestCommitHashAndDate(s.Project.Path)
	}


*/

func (s *Services) GetProjectInfoById(id string) (types.Project, error) {
	if id == "" {
		return types.Project{}, errors.New("id cannot be empty.")
	}

	if s.MainFolderInfo == nil {
		return types.Project{}, errors.New("scanFolderTypes cannot be empty.")
	}

	project, ok := s.MainFolderInfo.Index[id]
	if !ok {
		return types.Project{}, fmt.Errorf("No project with the id: %s", id)
	}

	return project, nil
}
