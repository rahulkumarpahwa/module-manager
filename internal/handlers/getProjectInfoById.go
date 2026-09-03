package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rahulkumarpahwa/module-manager/internal/utils"
	"github.com/rahulkumarpahwa/module-manager/services"
)

func (h *Handlers) GetProjectInfoById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "h.Services.GetProjectInfoById: ID can't be empty", http.StatusInternalServerError)
		return
	}

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


	projectDetail := services.ProjectDetails{
		ID:   id,
		Path: setPath.Path,
	}

	var projectDetails []services.ProjectDetails
	projectDetails = append(projectDetails, projectDetail)

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
