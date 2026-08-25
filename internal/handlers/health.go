package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Success     bool      `json:"success"`
	CurrentTime time.Time `json:"current_time"`
	Message     string    `json:"message"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Success:     true,
		CurrentTime: time.Now().Local(),
		Message:     "server is working fine.",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
