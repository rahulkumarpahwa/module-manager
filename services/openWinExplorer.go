package services

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (s *Services) OpenExplorer(projectPath string) error {
	if projectPath == "" {
		return fmt.Errorf("open explorer: project path can't be empty")
	}

	if runtime.GOOS != "windows" {
		return fmt.Errorf("open explorer: only supported on Windows")
	}

	cmd := exec.Command("explorer.exe", projectPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open explorer: %w", err)
	}

	return nil
}
