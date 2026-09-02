package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rahulkumarpahwa/module-manager/services/constants"
	"github.com/rahulkumarpahwa/module-manager/services/types"
)

func GenerateProjectDetails(projectPath string, entry os.DirEntry, i int) (*types.Project, error) {
	if !entry.IsDir() {
		return nil, errors.New("GenerateProjectDetails: entry can't be NON-DIR.")
	}

	if projectPath == "" {
		return nil, fmt.Errorf("GenerateProjectDetails: project path can't be empty.")
	}

	pe := types.Project{}

	pe.Name = entry.Name()
	pe.ID = generateID()
	pe.Path = projectPath

	var size int64
	size, err := GetDirSizeInBytes(pe.Path)
	if err != nil {
		return nil, fmt.Errorf("GenerateProjectDetails: %w", err)
	}
	pe.Size = size

	info, err := entry.Info()
	if err != nil {
		return nil, fmt.Errorf("GenerateProjectDetails: %w", err)
	}

	pe.ModifiedDate, pe.CreationDate = GetTimes(info)
	// .Format(time.RFC3339)

	entries, err := os.ReadDir(pe.Path)
	if err != nil {
		return nil, fmt.Errorf("GenerateProjectDetails: %w", err)
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

		switch name {
		case constants.NpmLockfile:
			pe.PackageManager = constants.NPM
		case constants.BunLockfile:
			pe.PackageManager = constants.BUN
		case constants.PnpmLockfile:
			pe.PackageManager = constants.PNPM
		case constants.YarnLockfile:
			pe.PackageManager = constants.YARN
		}

		underProject := types.UnderProjectDir{
			Name: name,
			// Size :  // fix this
			Path: filepath.Join(pe.Path, name),
		}

		pe.UnderDir = append(pe.UnderDir, underProject)
	}

	if pe.HasNodeModules {
		path := filepath.Join(pe.Path, "node_modules")
		size, err := GetNodeModulesSizeInBytes(path, constants.WORKERS)
		if err != nil {
			return nil, fmt.Errorf("GenerateProjectDetails: %w", err)

		}
		pe.NodeModulesSize = size
		pe.Size += pe.NodeModulesSize
	}

	if pe.HasGit {
		gitInfo, err := GetGitInfo(pe.Path)
		if err != nil {
			return nil, fmt.Errorf("GenerateProjectDetails : %w", err)
		}
		pe.LatestCommit = gitInfo.LatestCommit
		pe.LatestCommitDate = gitInfo.LatestCommitDate
	}

	return &pe, nil
}
