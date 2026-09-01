package handlers

import (
	"github.com/rahulkumarpahwa/module-manager/services"
	"github.com/rahulkumarpahwa/module-manager/services/types"
)

type Handlers struct {
	MainFolderInfo *types.MainFolder
	Services       *services.Services
}

func NewHandlers(mainFolderInfo *types.MainFolder,
	services *services.Services) *Handlers {
	return &Handlers{
		MainFolderInfo: mainFolderInfo,
		Services:       services,
	}
}
