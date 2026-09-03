package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rahulkumarpahwa/module-manager/internal/utils"
	"github.com/rahulkumarpahwa/module-manager/services"
)

func (h *Handlers) GetProjectInfos(w http.ResponseWriter, r *http.Request) {
	var projectDetails []services.ProjectDetails

	err := json.NewDecoder(r.Body).Decode(&projectDetails)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	projects, err := h.Services.GenerateProjectDetails(r.Context(), projectDetails)

	resp := utils.JsonResponse{
		Message: "project details by id",
		Status:  http.StatusOK,
		Success: true,
		Data: map[string]any{
			"projects": projects,
		},
	}
	resp.ResponseSend(w)

}
