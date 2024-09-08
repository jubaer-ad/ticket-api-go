package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

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

var collection *mongo.Collection
var client *mongo.Client

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	client, err = connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	collection = client.Database("TicketDB").Collection("tickets")

	app := fiber.New()
	app.Get("/health", health)

	app.Get("/api/Tickets", getTickets)
	// app.Get("/api/Tickets/:id", getTicketById)
	// app.Post("/api/Tickets", createTickets)
	// app.Patch("/api/Tickets/:id", updateTickets)
	// app.Delete("/api/Tickets/:id", deleteTickets)

	PORT := os.Getenv("PORT")

	fmt.Println("Server Starting")
	err = app.Listen(":" + PORT)
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnv() error {
	return godotenv.Load(".env")
}

func connectToDB() (*mongo.Client, error) {
	MONGODB_URL := os.Getenv("MONGODB_URL")

	clientOptions := options.Client().ApplyURI(MONGODB_URL)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer disconnect(client, ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func disconnect(client *mongo.Client, ctx context.Context) error {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func getCollection(dbName string, collectionName string) {

}

func health(c *fiber.Ctx) error {
	fmt.Printf("Health Check")
	return c.Status(200).JSON(fiber.Map{"msg": "Server is upNrunning!!!"})
}

func getTickets(c *fiber.Ctx) error {
	var tickets []Ticket
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(c.Context()) {
		var ticket Ticket
		if err := cursor.Decode(&ticket); err != nil {
			return err
		}
		tickets = append(tickets, ticket)
	}
	return c.JSON(tickets)
}
