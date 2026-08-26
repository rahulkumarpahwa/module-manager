package handlers

import (
	"module-manager/services"
	"module-manager/types"
)

type Handlers struct {
	MainFolderInfo *types.MainFolderInfo
	Services       *services.Services
}

func NewHandlers(mainFolderInfo *types.MainFolderInfo,
	services *services.Services) *Handlers {
	return &Handlers{
		MainFolderInfo: mainFolderInfo,
		Services:       services,
	}
}
