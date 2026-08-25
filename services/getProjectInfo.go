package services

import (
	"errors"
	"module-manager/types"
	"module-manager/utils"
	"os"
	"path/filepath"
)

func GetProjectInfo(sft *types.ScanFolderTypes, projectEntry *types.ProjectEntry) (*types.Project, error) {
	if sft == nil {
		return nil, errors.New("scanFolderTypes cannot be empty.")
	}
	if projectEntry == nil {
		return nil, errors.New("projectEntry cannot be empty.")
	}

	entries, err := os.ReadDir(projectEntry.Path)
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
		path := filepath.Join(projectEntry.Path, "node_modules")
		size, err := utils.GetDirSize(path)
		if err != nil {
			return nil, errors.New("GetProjectInfo: utils.GetDirSize: " + err.Error())
		}
		project.NodeModulesSize = size
	}

	if project.HasGit {
		project.LatestCommit, project.LatestCommitDate = utils.GetLatestCommitHashAndDate(projectEntry.Path)
	}

	return project, nil
}
