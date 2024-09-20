package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ticket-go/internal/database"
	"github.com/ticket-go/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// @Summary Server Health Check
// @Description Responds with a simple message for all method requests.
// @Tags health
// @Produce json
// @Success 200 {string} string "GET request accepted"
// @Failure 500 {string} string "Internal Server Error"
// @Router /health [get]
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

func TicketsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTicketsHandler(w, r)
	case http.MethodPost:
		createTicketHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// @Summary Get all tickets
// @Description Get all tickets
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {string} string "Get all tickets"
// @Router /api/tickets [get]
func getTicketsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := database.GetTickets()
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Create a new ticket
// @Description Creates a new ticket.
// @Tags tickets
// @Accept json
// @Produce json
// @Success 201 {string} string "Ticket created"
// @Router /api/tickets [post]
func createTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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
