package spots

import (
	"backend/models"
	"backend/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func Endpoints(app *fiber.App) {
	spots := app.Group("/spots")

	spots.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("PONG")
	})

	spots.Get("", func(c fiber.Ctx) error {
		spots, err := models.GetSpots()
		if err != nil {
			return utils.Error(c, err)
		}

		return c.JSON(spots)
	})

	spots.Get("/filter", func(c fiber.Ctx) error {
		// getting the paramters
		latitude, _ := strconv.ParseFloat(c.Query("latitude"), 32)
		longitude, _ := strconv.ParseFloat(c.Query("longitude"), 32)
		radius, _ := strconv.Atoi(c.Query("radius"))
		tags := c.Query("tags")

		center := []float64{longitude, latitude}

		splitTags := strings.Split(tags, ",")

		spots, err := models.GetSpotsFilter(center, radius, splitTags)
		if err != nil {
			return utils.Error(c, err)
		}

		return c.JSON(spots)
	})

	spots.Get("/search", func(c fiber.Ctx) error {
		input := c.Query("input")

		spots, err := models.SearchSpots(input)

		if err != nil {
			return utils.Error(c, err)
		}

		return c.JSON(spots)
	})
}
