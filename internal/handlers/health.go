package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Success     bool   `json:"success"`
	CurrentTime string `json:"current_time"`
	Message     string `json:"message"`
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Success:     true,
		CurrentTime: time.Now().In(time.FixedZone("IST", 5*60*60)).Format("2006-01-02 15:04:05"),
		Message: "server is working fine",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
