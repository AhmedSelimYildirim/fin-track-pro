package handlers

import (
	"fin-track-pro/internal/config"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())

	cfg := config.LoadConfig()

	userRepo := repository.NewUserRepository(repository.DB)
	marketRepo := repository.NewMarketRepository(repository.RedisClient)

	userService := service.NewUserService(userRepo)
	marketService := service.NewMarketService(cfg, marketRepo)

	userHandler := NewUserHandler(userService)
	marketHandler := NewMarketHandler(marketService)

	api := app.Group("/api")

	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Selam Ahmed Selim, sistem canavar gibi çalışıyor! 🚀",
		})
	})

	user := api.Group("/user")
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	market := api.Group("/market")
	market.Get("/rates", marketHandler.GetRates)

	assets := api.Group("/assets")
	assets.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Burada mal varlıkların listelenecek"})
	})
}
