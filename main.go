package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Ticket struct {
	title       string `json:"title"`
	description string `json:"description"`
	category    string `json:"category"`
	priority    int    `json:"priority"`
	progress    int    `json:"progress"`
	status      string `json:"status"`
	active      bool   `json:"active"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	MONGODB_URL := os.Getenv("MONGODB_URL")

	clientOptions := options.Client().ApplyURI(MONGODB_URL)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")

	collection = client.Database("TicketDB").Collection("tickets")

	app := fiber.New()
	app.Get("/health", health)

	PORT := os.Getenv("PORT")

	fmt.Println("Server Starting")
	err = app.Listen(":" + PORT)
	if err != nil {
		log.Fatal(err)
	}
}

func health(c *fiber.Ctx) error {
	fmt.Printf("Health Check")
	return c.Status(200).JSON(fiber.Map{"msg": "Server is upNrunning!!!"})
}
