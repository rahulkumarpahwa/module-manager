package utils

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type GitInfo struct {
	LatestCommit     string
	LatestCommitDate time.Time
}

// provides latest commit and
func GetGitInfo(projectPath string) (*GitInfo, error) {

	if projectPath == "" {
		return nil, fmt.Errorf("GetGitInfo: project path can't be empty.")
	}

	cmd := exec.Command("git", "-C", projectPath, "rev-parse", "HEAD")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	latestCommitDate, err := time.Parse(
		time.RFC3339,
		strings.TrimSpace(string(output)),
	)
	if err != nil {
		return nil, err
	}

	return &GitInfo{
		LatestCommit:     latestCommit,
		LatestCommitDate: latestCommitDate,
	}, nil

}
