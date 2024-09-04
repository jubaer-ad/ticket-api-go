package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Server Starting")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello from the Server"})
	})

	err := app.Listen(":4000")
	if err != nil {
		log.Fatal(err)
	}
}