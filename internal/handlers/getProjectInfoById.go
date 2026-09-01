package handlers

import (
	"net/http"

	"github.com/rahulkumarpahwa/module-manager/internal/utils"
)

func (h *Handlers) GetProjectInfoById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "h.Services.GetProjectInfoById: ID can't be empty", http.StatusInternalServerError)
		return
	}

	project, err := h.Services.GetEntryById(id)
	if err != nil {
		http.Error(w, "h.Services.GetProjectInfoById "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := utils.JsonResponse{
		Message: "project details by id",
		Status:  http.StatusOK,
		Success: true,
		Data: map[string]any{
			"project": project,
		},
	}
	resp.ResponseSend(w)

}
