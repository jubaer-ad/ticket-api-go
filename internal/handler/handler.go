package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ticket-go/internal/models"
)

func Health(w http.ResponseWriter, r *http.Request) {
	healthData := models.HealthResponse{
		Status:  "OK",
		Message: "Server is upNrunning!!!",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(healthData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
