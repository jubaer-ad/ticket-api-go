package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Server Starting")

	app := fiber.New()

	err := app.Listen(":4000")
	if err != nil {
		log.Fatal(err)
	}
}