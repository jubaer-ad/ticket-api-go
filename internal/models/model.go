package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Ticket represents a Ticket in the system.
// swagger:model Ticket
type Ticket struct {
	Id          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	Priority    int           `json:"priority"`
	Progress    int           `json:"progress"`
	Status      string        `json:"status"`
	Active      bool          `json:"active"`
	CreatedAt   *time.Time    `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   *time.Time    `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
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
