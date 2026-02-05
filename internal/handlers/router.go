package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())

	api := app.Group("/api")

	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Selam Ahmed Selim, sistem canavar gibi çalışıyor! 🚀",
		})
	})

	assets := api.Group("/assets")
	assets.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Burada mal varlıkların listelenecek"})
	})
}
