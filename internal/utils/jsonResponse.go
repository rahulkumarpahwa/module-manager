package utils

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Message string         `json:"message"`
	Status  int            `json:"status"`
	Success bool           `json:"success"`
	Data    map[string]any `json:"data"`
}

func (jr *JsonResponse) ResponseSend(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(jr.Status)
	json.NewEncoder(w).Encode(jr)
}
