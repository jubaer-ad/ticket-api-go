package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ticket-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func ConnectDB(cfg models.MongoDBConfig) {
	clientOptionsBuilder := options.Client().ApplyURI(cfg.MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(clientOptionsBuilder)
	if err != nil {
		log.Fatalf("Failed to connect to mongo: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	collection = client.Database(cfg.DatabaseName).Collection(cfg.CollectionName)
}

func GetTickets(ctx context.Context) (<-chan []bson.M, <-chan error) {
	ticketsChan := make(chan []bson.M, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(ticketsChan)
		defer close(errChan)

		var results []bson.M
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch tickets: %v", err)
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &results); err != nil {
			errChan <- fmt.Errorf("failed to decode tickets: %v", err)
			return
		}
		ticketsChan <- results
	}()

	return ticketsChan, errChan

}

func GetTicketByID(id bson.ObjectID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Ticket not found
		}
		return nil, err
	}

	return &ticket, nil
}

func CreateTicket(ticket models.Ticket) (*mongo.InsertOneResult, error) {
	rsp, err := collection.InsertOne(context.Background(), ticket)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func DeleteTicketByID(id bson.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter, options.Delete())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateTicketById(id bson.ObjectID, updateData models.Ticket) (*mongo.UpdateResult, error) {
	update := bson.M{
		"$set": bson.M{
			"title":       updateData.Title,
			"description": updateData.Description,
			"category":    updateData.Category,
			"Priority":    updateData.Priority,
			"Progress":    updateData.Progress,
			"status":      updateData.Status,
			"active":      updateData.Active,
			"updatedAt":   updateData.UpdatedAt,
		},
	}

	result, err := collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
