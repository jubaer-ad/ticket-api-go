package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Ticket struct {
	Id          bson.ObjectID `json:"id" bson:"_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	Priority    int           `json:"priority"`
	Progress    int           `json:"progress"`
	Status      string        `json:"status"`
	Active      bool          `json:"active"`
}
