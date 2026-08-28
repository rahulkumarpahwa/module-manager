package types

import "time"

type Project struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Path               string    `json:"path"`
	IsDir              bool      `json:"is_dir"`
	TotalSize          int64     `json:"total_size"`
	TotalSizeStr       string    `json:"total_size_str"`
	DateCreated        time.Time `json:"date_created"`
	DateModified       time.Time `json:"date_modified"`
	PackageManager     string   `json:"package_manager"`
	Libraries          []string  `json:"libraries"`
	HasPackageJson     bool      `json:"has_package_json"`
	HasNodeModules     bool      `json:"has_node_modules"`
	NodeModulesSize    int64     `json:"node_modules_size"`
	NodeModulesSizeStr string    `json:"node_modules_size_str"`
	HasGit             bool      `json:"has_git"`
	Frameworks         []string  `json:"frameworks"`
	LatestCommit       string    `json:"latest_commit"`
	LatestCommitDate   time.Time `json:"latest_commit_date"`
	UnderDir          []string  `json:"under_dir"`
}
