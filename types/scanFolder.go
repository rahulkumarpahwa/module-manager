package types

type ScanFolderTypes struct {
	Path         string `json:"path"`
	TotalSize    int64  `json:"total_size"`
	ProjectCount int    `json:"project_count"`
}
