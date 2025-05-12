package main

import (
	"backend/accounts"
	"backend/db"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	db.InitDB()
	db.InitCache()

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("PONG")
	})

	accounts.Endpoints(app)

	app.Listen(":3000")
}
