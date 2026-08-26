package handlers

import (
	"encoding/json"
	"module-manager/types"
	"module-manager/utils"
	"net/http"
	"strconv"
)

type ScanFolderResponse struct {
	Path          string               `json:"path"`
	TotalProjects int                  `json:"total_projects"`
	TotalSize     string               `json:"total_size"`
	Projects      []types.ProjectEntry `json:"projects"`
}

func (h *Handlers) ScanFolder(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.Services.ScanFolder(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp := &ScanFolderResponse{}

	resp.Path = h.MainFolderInfo.Path
	resp.TotalProjects = h.MainFolderInfo.ProjectCount
	resp.TotalSize = utils.FormatSize(h.MainFolderInfo.TotalSize)
	resp.Projects = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
