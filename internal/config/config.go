package config

import (
	"os"

	"github.com/ticket-go/internal/models"
)

func LoadMongoConfig() models.MongoDBConfig {
	return models.MongoDBConfig{
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
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
