package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ticket-go-fiber/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var collection *mongo.Collection

func LoadCollection(c *mongo.Collection) {
	collection = c
}

func Health(c *fiber.Ctx) error {
	fmt.Println("Health Check")
	return c.Status(200).JSON(fiber.Map{"msg": "Server is up and running!!!"})
}

func GetTickets(c *fiber.Ctx) error {
	var tickets []models.Ticket
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(c.Context()) {
		var ticket models.Ticket
		if err := cursor.Decode(&ticket); err != nil {
			return err
		}
		tickets = append(tickets, ticket)
	}
	return c.JSON(tickets)
}
