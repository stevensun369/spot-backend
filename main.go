package main

import (
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("PONG")
	})

	app.Listen(":3000")
}
