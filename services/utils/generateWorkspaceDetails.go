package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rahulkumarpahwa/module-manager/services/types"
)

func GenerateWorkspaceDetails(entryPath string, entry os.DirEntry, i int) (*types.WorkSpace, error) {
	ws := types.WorkSpace{}
	ws.SNo = i
	ws.ID = generateID()
	ws.Name = entry.Name()

	info, err := os.Stat(entryPath)
	if err != nil {
		return nil, fmt.Errorf("GenerateWorkspaceDetails: %w", err)
	}
	ws.CreationDate, ws.ModifiedDate = GetTimes(info)

	var wsp []types.Project
	var wg sync.WaitGroup

	err = filepath.WalkDir(entryPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		// ws.UnderWorkSpace = append(ws.UnderWorkSpace, info.Name())

		projectPath := filepath.Join(entryPath, entry.Name())

		if getMatchCases(entry.Name()) {
			ws.UnderWorkSpace = append(ws.UnderWorkSpace, info.Name())

			wg.Add(1)

			go func(projectPath string, entry fs.DirEntry, i int) (*types.Project, error) {

				defer wg.Done()

				pe, err := GenerateProjectDetails(projectPath, entry, i)
				if err != nil {
					return nil, fmt.Errorf("GenerateWorkspaceDetails : %w", err)
				}

				wsp = append(wsp, *pe)
				return nil, nil
			}(projectPath, entry, 1)

		}

		if !getMatchCases(entry.Name()) {
			ws.Size += info.Size()
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("GenerateWorkspaceDetails : %w", err)
	}

	wg.Wait()

	ws.WorkSpaceProjects = append(ws.WorkSpaceProjects, wsp...)
	ws.WorkSpaceProjectsCount = len(wsp)

	for _, p := range wsp {
		ws.Size += p.Size
	}

	return &ws, nil
}

func getMatchCases(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))

	switch name {
	case "frontend", "front-end", "front_end":
		return true

	case "backend", "back-end", "back_end":
		return true

	default:
		return false
	}
}
