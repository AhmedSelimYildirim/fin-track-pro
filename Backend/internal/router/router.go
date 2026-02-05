package router

import (
	"fin-track-pro/internal/config"
	"fin-track-pro/internal/handlers"
	"fin-track-pro/internal/middleware"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	cfg := config.LoadConfig()

	userRepo := repository.NewUserRepository(repository.DB)
	assetRepo := repository.NewAssetRepository(repository.DB)
	calendarRepo := repository.NewCalendarRepository(repository.DB)
	marketRepo := repository.NewMarketRepository(repository.DB)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	marketService := service.NewMarketService(cfg, repository.RedisClient)
	assetService := service.NewAssetService(assetRepo, marketService)
	calendarService := service.NewCalendarService(calendarRepo)

	userHandler := handlers.NewUserHandler(userService)
	marketHandler := handlers.NewMarketHandler(marketService)
	assetHandler := handlers.NewAssetHandler(assetService, marketService)
	calendarHandler := handlers.NewCalendarHandler(calendarService)

	_ = marketRepo

	api := app.Group("/api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Sistem hazır ! 🚀"})
	})

	user := api.Group("/user")
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	market := api.Group("/market")
	market.Get("/rates", marketHandler.GetRates)

	protected := api.Group("/", middleware.Protected(cfg.JWTSecret))

	u := protected.Group("/user")
	u.Put("/update", userHandler.Update)
	u.Delete("/delete", userHandler.Delete)

	assets := protected.Group("/assets")
	assets.Post("/update", assetHandler.UpdateBalance)
	assets.Get("/summary", assetHandler.GetSummary)
	assets.Get("/transactions", assetHandler.GetTransactions)
	assets.Get("/receipt/:id", assetHandler.GetReceipt)

	calendar := protected.Group("/calendar")
	calendar.Post("/remind", calendarHandler.AddEvent)
	calendar.Get("/list", calendarHandler.ListReminders)
}
