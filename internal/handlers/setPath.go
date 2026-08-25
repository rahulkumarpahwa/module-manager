package handlers

import (
	"encoding/json"
	"module-manager/types"
	"net/http"
)

type SetPathType struct {
	Path string `json:"path"`
}

func SetPath(scan *types.ScanFolderTypes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		scan.Path = setPath.Path

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(struct {
			success string
		}{
			success: "ok",
		})
	}
}
