package types

type MainFolderInfo struct {
	Path         string             `json:"path"`
	TotalSize    int64              `json:"total_Size"`
	TotalSizeStr string             `json:"total_size_str"`
	ProjectCount int                `json:"project_count"`
	Projects     []Project          `json:"projects"`
	Index        map[string]Project `json:"-"`
}

func NewMainFolderInfo() *MainFolderInfo {
	return &MainFolderInfo{}
}
