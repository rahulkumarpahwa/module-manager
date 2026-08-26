package types

type MainFolderInfo struct {
	Path         string
	TotalSize    int64
	ProjectCount int
	Index        map[int64]Project
}

func NewMainFolderInfo() *MainFolderInfo {
	return &MainFolderInfo{}
}
