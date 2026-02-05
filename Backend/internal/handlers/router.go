package handlers

import (
	"fin-track-pro/internal/config"
	"fin-track-pro/internal/middleware"
	repository2 "fin-track-pro/internal/repository"
	service2 "fin-track-pro/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	cfg := config.LoadConfig()

	userRepo := repository2.NewUserRepository(repository2.DB)
	assetRepo := repository2.NewAssetRepository(repository2.DB)
	calendarRepo := repository2.NewCalendarRepository(repository2.DB)

	userService := service2.NewUserService(userRepo, cfg.JWTSecret)
	marketService := service2.NewMarketService(cfg, repository2.RedisClient)
	assetService := service2.NewAssetService(assetRepo, marketService)
	calendarService := service2.NewCalendarService(calendarRepo)

	userHandler := NewUserHandler(userService)
	marketHandler := NewMarketHandler(marketService)
	assetHandler := NewAssetHandler(assetService, marketService)
	calendarHandler := NewCalendarHandler(calendarService)

	api := app.Group("/api")
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Sistem hazır Ahmed Selim! 🚀"})
	})

	user := api.Group("/user")
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	market := api.Group("/market")
	market.Get("/rates", marketHandler.GetRates)

	protected := api.Group("/", middleware.Protected(cfg.JWTSecret))

	assets := protected.Group("assets")
	assets.Post("/", assetHandler.CreateAsset)
	assets.Get("/", assetHandler.GetAssets)
	assets.Get("/summary", assetHandler.GetSummary)
	assets.Get("/export", assetHandler.ExportExcel)
	assets.Get("/receipt/:id", assetHandler.GetReceipt)
	assets.Put("/:id", assetHandler.UpdateAsset)
	assets.Delete("/:id", assetHandler.DeleteAsset)

	calendar := protected.Group("calendar")
	calendar.Post("/remind", calendarHandler.AddEvent)
}
