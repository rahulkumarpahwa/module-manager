package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rahulkumarpahwa/module-manager/services/types"
)

func GenerateWorkspaceDetails(
	entryPath string,
	entry os.DirEntry,
	i int,
) (*types.WorkSpace, error) {

	ws := types.WorkSpace{
		SNo:  i,
		ID:   GenerateID(),
		Name: entry.Name(),
	}

	info, err := entry.Info()
	if err != nil {
		return nil, fmt.Errorf("generate workspace details: %w", err)
	}

	ws.CreationDate, ws.ModifiedDate = GetTimes(info)

	entries, err := os.ReadDir(entryPath)
	if err != nil {
		return nil, fmt.Errorf("read workspace directory: %w", err)
	}

	var (
		wg      sync.WaitGroup
		wspChan = make(chan types.Project)
	)

	for i, projectEntry := range entries {
		if !projectEntry.IsDir() {
			ws.Size += getEntrySize(filepath.Join(entryPath, projectEntry.Name()))
			continue
		}

		projectName := projectEntry.Name()

		if !getMatchCases(projectName) {
			projectInfo, err := projectEntry.Info()
			if err != nil {
				return nil, fmt.Errorf("get entry info: %w", err)
			}

			ws.Size += projectInfo.Size()
			continue
		}

		ws.UnderWorkSpace = append(
			ws.UnderWorkSpace,
			projectName,
		)

		projectPath := filepath.Join(entryPath, projectName)

		wg.Add(1)

		go func(
			projectPath string,
			projectEntry os.DirEntry,
		) {
			defer wg.Done()

			project, err := GenerateProjectDetails(
				projectPath,
				projectEntry,
				i,
			)
			if err != nil {
				return
			}

			wspChan <- *project

		}(projectPath, projectEntry)
	}

	go func() {
		wg.Wait()
		close(wspChan)
	}()

	wsp := []types.Project{}

	for result := range wspChan {
		wsp = append(wsp, result)
	}

	ws.WorkSpaceProjects = append(ws.WorkSpaceProjects, wsp...)
	ws.WorkSpaceProjectsCount = len(wsp)

	for _, project := range wsp {
		ws.Size += project.Size
	}

	return &ws, nil
}

func getMatchCases(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))

	switch name {
	case "frontend", "front-end", "front_end",
		"backend", "back-end", "back_end":
		return true
	default:
		return false
	}
}

func getEntrySize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return info.Size()
}
