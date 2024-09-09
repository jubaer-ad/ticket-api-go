package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ticket-go-fiber/handlers"
)

func LoadRoutes(app *fiber.App) {

	app.Get("/health", handlers.Health)

	app.Get("/api/Tickets", handlers.GetTickets)
	// app.Get("/api/Tickets/:id", getTicketById)
	// app.Post("/api/Tickets", createTickets)
	// app.Patch("/api/Tickets/:id", updateTickets)
	// app.Delete("/api/Tickets/:id", deleteTickets)
}
