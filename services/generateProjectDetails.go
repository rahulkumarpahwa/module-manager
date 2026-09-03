package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rahulkumarpahwa/module-manager/services/constants"
	"github.com/rahulkumarpahwa/module-manager/services/types"
	"github.com/rahulkumarpahwa/module-manager/services/utils"
)

type ProjectDetails struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type Results struct {
	project *types.Project
	err     error
}

func (s *Services) GenerateProjectDetails(ctx context.Context, projectDetails []ProjectDetails) ([]types.Project, error) {
	jobs := make(chan ProjectDetails)
	results := make(chan Results, len(projectDetails))

	const workers = 5
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-jobs:
					if !ok {
						return
					}

					pe, err := generateProjectDetailsWorker(ctx, item.Path, item.ID)

					select {
					case results <- Results{
						project: pe,
						err:     err,
					}:
					case <-ctx.Done():
						return
					}

				}

			}

		}()
	}

	go func() {
		defer close(jobs)

		for _, item := range projectDetails {
			select {
			case jobs <- item:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var projects []types.Project

	for result := range results {
		if result.err != nil {
			// If cancellation happened, return it.
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			return nil, result.err
		}

		if result.project != nil {
			projects = append(projects, *result.project)
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func generateProjectDetailsWorker(ctx context.Context, path string, ID string) (*types.Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if path == "" {
		return nil, fmt.Errorf("GenerateProjectDetails: project path can't be empty.")
	}

	pe := types.Project{
		ID:   ID,
		Path: path,
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var size int64
	size, err := utils.GetDirSizeInBytes(ctx, pe.Path)
	if err != nil {
		return nil, fmt.Errorf("GenerateProjectDetails: %w", err)
	}
	pe.Size = size

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(pe.Path)
	if err != nil {
		return nil, fmt.Errorf("GenerateProjectDetails: %w", err)
	}

	for _, entry := range entries {

		if err := ctx.Err(); err != nil {
			return nil, err
		}

		var name = entry.Name()

		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("GenerateProjectDetails: %w", err)
		}

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
			Size: info.Size(),
			Path: filepath.Join(pe.Path, name),
		}

		pe.UnderDir = append(pe.UnderDir, underProject)
	}

	if pe.HasNodeModules {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		path := filepath.Join(pe.Path, "node_modules")
		size, err := utils.GetNodeModulesSizeInBytes(ctx, path, constants.WORKERS)
		if err != nil {
			return nil, fmt.Errorf("GenerateProjectDetails: %w", err)

		}
		pe.NodeModulesSize = size
		pe.Size += pe.NodeModulesSize
	}

	if pe.HasGit {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		gitInfo, err := utils.GetGitInfo(pe.Path)
		if err != nil {
			return nil, fmt.Errorf("GenerateProjectDetails : %w", err)
		}
		pe.LatestCommit = gitInfo.LatestCommit
		pe.LatestCommitDate = gitInfo.LatestCommitDate
	}

	return &pe, nil
}

/*
pesChan := make(chan types.Project)
	wsChan := make(chan types.WorkSpace)

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

	s.MainFolder.IndexProject = makeIndexProject(s.MainFolder.Projects)
	s.MainFolder.IndexWorkspace = makeIndexWorkspace(s.MainFolder.WorkSpaces)

*/
