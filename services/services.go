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

/**
usersByID := make(map[int64]User, len(users))

for _, user := range users {
	usersByID[user.ID] = user
}

user, ok := usersByID[42]
if !ok {
	// not found
	return
}

fmt.Println(user.Name)

*/

func (s *Services) GetProjectInfoById() (*types.Project, error) {
	if s.MainFolderInfo == nil {
		return nil, errors.New("scanFolderTypes cannot be empty.")
	}
	if s.Project == nil {
		return nil, errors.New("projectEntry cannot be empty.")
	}

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

	return project, nil
}
