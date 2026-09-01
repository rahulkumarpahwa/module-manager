package types

import "time"

type MainFolder struct {
	Path            string         `json:"path"`
	TotalSize       int64          `json:"total_Size"`
	CreationDate    time.Time      `json:"creation_date"`
	ModifiedDate    time.Time      `json:"modified_date"`
	ProjectsCount   int            `json:"projects_count"`
	WorkSpacesCount int            `json:"workspaces_count"`
	Projects        []Project      `json:"projects"`
	WorkSpaces      []WorkSpace    `json:"workspaces"`
	Index           map[string]any `json:"-"`
}

func NewMainFolder() *MainFolder {
	return &MainFolder{}
}
