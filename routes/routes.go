package routes

import (
	"net/http"

	"github.com/ticket-go/internal/handler"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	return mux
}
