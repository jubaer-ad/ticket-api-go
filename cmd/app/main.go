package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/ticket-go/docs"
	"github.com/ticket-go/internal/config"
	"github.com/ticket-go/internal/database"
	"github.com/ticket-go/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cfg := config.LoadMongoConfig()
	database.ConnectDB(cfg)
	router := routes.NewRouter()
	port := os.Getenv("PORT")
	url := os.Getenv("SERVER_URL")
	addr := fmt.Sprintf("%s:%s", url, port)
	fmt.Println(addr, port, url)
	fmt.Printf("Server is listening on %s\n", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
