package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ticket-go/internal/database"
	"github.com/ticket-go/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	var dbRsp *mongo.InsertOneResult
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dbRsp, err = database.CreateTicket(ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rsp := struct {
		ID primitive.ObjectID `josn:"id"`
	}{
		ID: dbRsp.InsertedID.(primitive.ObjectID),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rsp)
}
