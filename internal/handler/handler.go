package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ticket-go/internal/database"
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

func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := database.GetTickets()
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
	}
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling data", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
