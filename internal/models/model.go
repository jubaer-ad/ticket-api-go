package models

// Ticket represents a Ticket in the system.
// swagger:model Ticket
type Ticket struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Priority    int    `json:"priority"`
	Progress    int    `json:"progress"`
	Status      string `json:"status"`
	Active      bool   `json:"active"`
}

// HealthResponse represents thr health check response.
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
