package database

import (
	"context"
	"log"

	"github.com/ticket-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func ConnectDB(cfg models.MongoDBConfig) {
	var err error
	clientOptionsBuilder := options.Client().ApplyURI(cfg.MongoURI)
	client, err = mongo.Connect(clientOptionsBuilder)
	if err != nil {
		log.Fatalf("Failed to connect to mongo: %v", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	collection = client.Database(cfg.DatabaseName).Collection(cfg.CollectionName)
}

func GetTickets() ([]bson.M, error) {
	var results []bson.M
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil

}

func CreateTicket(ticket models.Ticket) (*mongo.InsertOneResult, error) {
	rsp, err := collection.InsertOne(context.Background(), ticket)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
