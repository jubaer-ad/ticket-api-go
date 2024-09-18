package config

import (
	"os"

	"github.com/ticket-go/internal/models"
)

func LoadMongoConfig() models.MongoDBConfig {
	return models.MongoDBConfig{
		MongoURI:       getEnv("MONGO_URI", "mongodb+srv://admin:Awb3WeJjFvMOKpMB@cluster0.b2v96.mongodb.net/TicketDB"),
		DatabaseName:   getEnv("DATABASE_NAME", "TicketDB"),
		CollectionName: getEnv("COLLECTION_NAME", "tickets"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
