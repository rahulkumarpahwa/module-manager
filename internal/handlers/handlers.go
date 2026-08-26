package handlers

import (
	"github.com/rahulkumarpahwa/module-manager/services"
	"github.com/rahulkumarpahwa/module-manager/types"
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
