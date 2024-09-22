package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ticket-go/internal/database"
	"github.com/ticket-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// @Summary Server Health Check
// @Description Responds with a simple message for all method requests.
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse "Ok"
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
		if r.URL.Path == "/api/tickets" {
			getTicketsHandler(w, r)
		} else {
			getTicketHandler(w, r)
		}
	case http.MethodPost:
		createTicketHandler(w, r)
	case http.MethodDelete:
		deleteTicketHandler(w, r)
	case http.MethodPut, http.MethodPatch:
		updateTicketHandler(w, r)
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
// @Success 200 {array} models.Ticket "Ok"
// @Router /api/tickets [get]
func getTicketsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ticketsChan, errChan := database.GetTickets(ctx)
	var once sync.Once

	go func() {
		select {
		case tickets := <-ticketsChan:
			once.Do(func() {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(tickets); err != nil {
					http.Error(w, "Error encoding response", http.StatusInternalServerError)
				}
			})
		case err := <-errChan:
			once.Do(func() {
				if err != nil {
					http.Error(w, "Error fetching data", http.StatusInternalServerError)
				}
			})
		case <-ctx.Done():
			once.Do(func() {
				http.Error(w, "Request timed out", http.StatusRequestTimeout)
			})
		}
	}()
}

// @Summary Get a ticket by ID
// @Description Get a ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {object} models.Ticket "Success"
// @Failure 404 {string} string "Ticket not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tickets/{id} [get]
func getTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/tickets/"):]
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := database.GetTicketByID(objectID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

// @Summary Create a new ticket
// @Description Creates a new ticket.
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body models.Ticket true "Ticket object"
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
	now := time.Now()
	ticket.CreatedAt = &now

	dbRsp, err = database.CreateTicket(ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rsp := struct {
		ID bson.ObjectID `josn:"id"`
	}{
		ID: dbRsp.InsertedID.(bson.ObjectID),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rsp)
}

// @Summary Delete a ticket by ID
// @Description Deletes a ticket from the MongoDB collection by BSON ID
// @Tags tickets
// @Produce json
// @Param id path string true "Ticket ID"
// @Success 200 {string} string "Ticket deleted"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Ticket not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tickets/{id} [delete]
func deleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Path[len("/api/tickets/"):]
	objID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DeleteTicketByID(objID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rsp := models.SuccessResponse{
		Message: "Ticket deleted successfully",
	}
	json.NewEncoder(w).Encode(rsp)
}

// @Summary Update a ticket by ID
// @Description Updates a ticket in the MongoDB collection by its ID.
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID"
// @Param ticket body models.Ticket true "Updated Ticket object"
// @Success 200 {string} string "Ticket updated"
// @Failure 400 {string} string "Invalid ID format or bad request"
// @Failure 404 {string} string "Ticket not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/tickets/{id} [put]
func updateTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/tickets/"):]
	objID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateData models.Ticket
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	updateData.UpdatedAt = &now

	result, err := database.UpdateTicketById(objID, updateData)
	if err != nil {
		http.Error(w, "Error updating ticket", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rsp := models.SuccessResponse{
		Message: "Ticket updated successfully",
	}
	json.NewEncoder(w).Encode(rsp)
}
