package handlers

import (
	"encoding/json"
	"net/http"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Success string `json:"success"`
	}{
		Success: "ok",
	})
}
