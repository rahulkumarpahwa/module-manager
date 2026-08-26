package types

type MainFolderInfo struct {
	Path         string
	TotalSize    int64
	ProjectCount int
	Projects     []Project
	Index        map[string]Project
}

func NewMainFolderInfo() *MainFolderInfo {
	return &MainFolderInfo{}
}
