package types

import "time"

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`

	CreationDate time.Time `json:"date_created"`
	ModifiedDate time.Time `json:"date_modified"`

	PackageManager string `json:"package_manager"`

	HasPackageJson  bool  `json:"has_package_json,omitempty"` // currently we are working with npm and later we will work with other pakage managers.
	HasNodeModules  bool  `json:"has_node_modules,omitempty"`
	NodeModulesSize int64 `json:"node_modules_size,omitempty"`

	HasGit           bool      `json:"git_status,omitempty"`
	LatestCommit     string    `json:"latest_commit,omitempty"`
	LatestCommitDate time.Time `json:"latest_commit_date,omitempty"`

	Libraries []string `json:"libraries,omitempty"`

	UnderDir []UnderProjectDir `json:"under_dir,omitempty"`
}

type UnderProjectDir struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
	Path string `json:"path,omitempty"`
}

func (p Project) GetID() string {
	return p.ID
}
