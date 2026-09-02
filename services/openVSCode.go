package services

import (
	"fmt"
	"os/exec"
)

func (s *Services) OpenVSCode(projectPath string) error {
	if projectPath == "" {
		return fmt.Errorf("open vscode: project path can't be empty")
	}

	cmd := exec.Command("code", projectPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("open vscode: %w", err)
	}

	return nil
}

/**
One important thing

This requires the VS Code code command to be available in your PATH.

You can test that from Git Bash:

code --version

If you get something like:

1.XX.X

your Go code should be able to execute code.

If you get:

bash: code: command not found

then VS Code's command-line launcher isn't in your PATH.
*/
