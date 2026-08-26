package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rahulkumarpahwa/module-manager/internal/utils"
)

type SetPathType struct {
	Path string `json:"path"`
}

func (h *Handlers) SetPath(w http.ResponseWriter, r *http.Request) {
	var setPath SetPathType

	err := json.NewDecoder(r.Body).Decode(&setPath)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if setPath.Path == "" {
		http.Error(w, "path cannot be empty", http.StatusBadRequest)
		return
	}

	h.MainFolderInfo.Path = setPath.Path

	resp := utils.JsonResponse{
		Message: "path is set",
		Status:  http.StatusOK,
		Success: true,
		Data: map[string]any{
			"path" : h.MainFolderInfo.Path,
		},
	}
	resp.ResponseSend(w)
}
