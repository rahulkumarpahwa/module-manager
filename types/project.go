package types

import "time"

type Project struct {
	ProjectEntry
	PackageManager   string    `json:"package_manager"`
	Libraries        []string  `json:"libraries"`
	HasPackageJson   bool      `json:"has_package_json"`
	HasNodeModules   bool      `json:"has_node_modules"`
	NodeModulesSize  int64     `json:"node_modules_size"`
	HasGit           bool      `json:"has_git"`
	Frameworks       []string  `json:"frameworks"`
	LatestCommit     string    `json:"latest_commit"`
	LatestCommitDate time.Time `json:"latest_commit_date"`
	ListFiles        []string  `json:"list_files"`
}
