package accounts

import (
	"backend/models"

	"github.com/gofiber/fiber/v3"
)

func Endpoints(app *fiber.App) {
	acc := app.Group("/accounts")

	acc.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("PONG")
	})

	onboarding := acc.Group("/onboarding")
	onboarding.Post("", postOnboarding)
	onboarding.Post("/verify-code", postVerifyCode, models.AccountMiddleware)
	onboarding.Post("/name", postName, models.AccountMiddleware)

	acc.Post("/bio", postBio, models.AccountMiddleware)
	acc.Patch("/socials", patchSocial, models.AccountMiddleware)

	acc.Post("/follow", postFollow, models.AccountMiddleware)
	acc.Post("/unfollow", postUnfollow, models.AccountMiddleware)
}
