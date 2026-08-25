package utils

import (
	"os/exec"
	"strings"
	"time"
)

func GetLatestCommitHashAndDate(projectPath string) (string, time.Time) {
	cmd := exec.Command("git", "-C", projectPath, "rev-parse", "HEAD")

	output, err := cmd.Output()
	if err != nil {
		// Not a Git repository or Git isn't available
		return "", time.Time{}
	}

	latestCommit := strings.TrimSpace(string(output))

	cmd = exec.Command(
		"git",
		"-C", projectPath,
		"log",
		"-1",
		"--format=%cI",
	)

	output, err = cmd.Output()
	if err != nil {
		return "", time.Time{}
	}

	latestCommitDate, err := time.Parse(
		time.RFC3339,
		strings.TrimSpace(string(output)),
	)
	if err != nil {
		return "", time.Time{}
	}

	return latestCommit, latestCommitDate
}
