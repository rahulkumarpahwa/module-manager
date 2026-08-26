package handlers

import (
	"encoding/json"
	"github.com/rahulkumarpahwa/module-manager/types"
	"github.com/rahulkumarpahwa/module-manager/utils"
	"net/http"
	"strconv"
)

type ScanMainFolderResponse struct {
	Path          string          `json:"path"`
	TotalProjects int             `json:"total_projects"`
	TotalSize     string          `json:"total_size"`
	Projects      []types.Project `json:"projects"`
}

func (h *Handlers) ScanMainFolder(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit := 20 // default
	offset := 0 // default

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		limit = parsed
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		offset = parsed
	}

	data, err := h.Services.ScanMainFolder(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp := &ScanMainFolderResponse{}

	resp.Path = h.MainFolderInfo.Path
	resp.TotalProjects = h.MainFolderInfo.ProjectCount
	resp.TotalSize = utils.FormatSize(h.MainFolderInfo.TotalSize)
	resp.Projects = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
