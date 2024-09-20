package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ticket represents a Ticket in the system.
// swagger:model Ticket
type Ticket struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Priority    int                `json:"priority"`
	Progress    int                `json:"progress"`
	Status      string             `json:"status"`
	Active      bool               `json:"active"`
}

// HealthResponse represents the health check response.
// swagger:model HealthResponse
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MongoDBConfig struct {
	MongoURI       string
	DatabaseName   string
	CollectionName string
}

// SuccessResponse represents a success response.
// swagger:model SuccessResponse
type SuccessResponse struct {
	Message string `json:"message"`
}
