package types

import (
	"time"
)

type WorkSpace struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	WorkSpaceProjects      []Project `json:"projects"`
	WorkSpaceProjectsCount int             `json:"projects_count"`
	CreationDate          time.Time       `json:"creation_date"`
	ModifiedDate          time.Time       `json:"modified_date"`
	UnderWorkSpace        []string        `json:"under_workspace,omitempty"`
	Size                  int64           `json:"size"`
}

// workspace contains the many project and whereas the project is one where we have the node_modules as the workspace can have the package.json in case of the multiple projects in it with different node_modules of each. we consider the project one who has the node_modules.
