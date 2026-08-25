package handlers

import (
	"encoding/json"
	"module-manager/services"
	"module-manager/types"
	"net/http"
)

type ScanFolderResponse struct {
	Path          string               `json:"path"`
	TotalProjects int                  `json:"total_projects"`
	TotalSize     int64                `json:"total_size"`
	Projects      []types.ProjectEntry `json:"projects"`
}

func ScanFolder(scan *types.ScanFolderTypes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := services.ScanFolder(scan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		resp := &ScanFolderResponse{}

		resp.Path = scan.Path
		resp.TotalProjects = scan.ProjectCount
		resp.TotalSize = scan.TotalSize
		resp.Projects = data

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
