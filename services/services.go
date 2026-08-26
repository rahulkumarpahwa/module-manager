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
}

func NewServices(mainFolderInfo *types.MainFolderInfo) *Services {
	return &Services{
		MainFolderInfo: mainFolderInfo,
	}
}

func (s *Services) ScanMainFolder(limit int, offset int) (*types.MainFolderInfo, error) {
	if s.MainFolderInfo.Path == "" {
		return nil, errors.New("path must exist")
	}

	path := filepath.Clean(s.MainFolderInfo.Path)

	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("folder does not exist")
		return nil, errors.New("folder does not exist")
	}

	if !info.IsDir() {
		log.Fatalf("path is not a directory/folder")
		return nil, errors.New("path is not a directory/folder")
	}

	projectEntries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("os.ReadDir: " + err.Error())
		return nil, errors.New("os.ReadDir: " + err.Error())
	}

	var pes []types.Project

	entries := take(projectEntries, offset, limit)

	s.MainFolderInfo.TotalSize = 0

	for i, entry := range entries {

		if !entry.IsDir() {
			// when entry is not a Dir we will not consider.
			continue
		}
		pe, err := getProjectDetails(path, entry, i)
		if err != nil {
			return nil, errors.New("ScanMainFolder.getProjectDetails: " + err.Error())
		}

		s.MainFolderInfo.TotalSize += pe.TotalSize
		pes = append(pes, pe)
	}

	s.MainFolderInfo.ProjectCount = len(entries)
	s.MainFolderInfo.Projects = pes
	s.MainFolderInfo.Index = makeIndex(pes)

	return s.MainFolderInfo, nil
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

func getProjectDetails(path string, entry os.DirEntry, i int) (types.Project, error) {
	if !entry.IsDir() {
		return types.Project{}, errors.New("entry can't be non dir.")
	}

	pe := types.Project{}

	pe.IsDir = true
	pe.Name = entry.Name()
	pe.ID = fmt.Sprintf("%d%s", i+1, pe.Name)
	pe.Path = filepath.Join(path, pe.Name)

	var size utils.DirSize
	size, err := utils.GetDirSize(pe.Path)
	if err != nil {
		log.Fatal("utils.GetDirSize: " + err.Error())
		return types.Project{}, errors.New("utils.GetDirSize: " + err.Error())
	}
	pe.TotalSize = size.Bytes

	entries, err := os.ReadDir(pe.Path)
	if err != nil {
		return types.Project{}, errors.New("services.getProjectDetails: os.ReadDir: " + err.Error())
	}

	for _, entry := range entries {
		var name = entry.Name()

		if name == "package.json" {
			pe.HasPackageJson = true
		}

		if name == "node_modules" {
			pe.HasNodeModules = true
		}

		if name == ".git" {
			pe.HasGit = true
		}

		pe.ListFiles = append(pe.ListFiles, name)
	}

	if pe.HasNodeModules {
		path := filepath.Join(pe.Path, "node_modules")
		size, err := utils.GetDirSize(path)
		if err != nil {
			return types.Project{}, errors.New("services.getProjectDetails: utils.GetDirSize: " + err.Error())

		}
		pe.NodeModulesSize = size.Bytes
	}

	if pe.HasGit {
		pe.LatestCommit, pe.LatestCommitDate = utils.GetLatestCommitHashAndDate(pe.Path)
	}

	return pe, nil
}

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
