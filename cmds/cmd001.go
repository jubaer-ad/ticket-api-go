package cmds

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/ticket-go-fiber/handlers"
	"github.com/ticket-go-fiber/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var collection *mongo.Collection
var client *mongo.Client
var app *fiber.App

func getClient() *mongo.Client {
	return client
}

func loadEnv() error {
	return godotenv.Load(".env")
}

func connectToDB() error {
	MONGODB_URL := os.Getenv("MONGODB_URL")

	clientOptions := options.Client().ApplyURI(MONGODB_URL)
	var err error
	client, err = mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer disconnect(client, ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

func disconnect(client *mongo.Client, ctx context.Context) error {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func GetCollection() *mongo.Collection {
	return collection
}

func loadCollection(dbName string, collectionName string) {
	collection = client.Database(dbName).Collection(collectionName)
	handlers.LoadCollection(collection)
}

func LoadCMD001(dbName string, collectionName string) {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	loadCollection(dbName, collectionName)

	app := fiber.New()
	routes.LoadRoutes(app)
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
