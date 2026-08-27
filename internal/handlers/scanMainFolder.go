package handlers

import (
	"net/http"
	"strconv"

	"github.com/rahulkumarpahwa/module-manager/internal/utils"
)


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
		return
	}


	resp := utils.JsonResponse{
		Message: "Main Path Scan Completed",
		Status:  http.StatusOK,
		Success: true,
		Data: map[string]any{
			"directory_details": data,
		},
	}
	resp.ResponseSend(w)
}
