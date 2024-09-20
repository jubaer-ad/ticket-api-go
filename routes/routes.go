package routes

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/ticket-go/internal/handler"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	// Swagger documentation route
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/health", handler.Health)

	mux.HandleFunc("/api/tickets", handler.TicketsHandler)
	mux.HandleFunc("/api/tickets/{id}", handler.TicketsHandler)

	return mux
}
