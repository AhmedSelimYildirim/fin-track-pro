package router

import (
	"fin-track-pro/internal/handler"
	"fin-track-pro/internal/infrastructure/config"
	"fin-track-pro/internal/infrastructure/database"
	"fin-track-pro/internal/infrastructure/redis"
	"fin-track-pro/internal/middleware"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://fin-track-pro-1.onrender.com, https://fin-track-pro.onrender.com",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Currency, X-Ayar",
		AllowCredentials: true,
	}))

	cfg := config.LoadConfig()

	userRepo := repository.NewUserRepository(database.DB)
	assetRepo := repository.NewAssetRepository(database.DB)
	calendarRepo := repository.NewCalendarRepository(database.DB)
	marketRepo := repository.NewMarketRepository(database.DB)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	marketService := service.NewMarketService(cfg, redis.Client, marketRepo)
	assetService := service.NewAssetService(assetRepo, marketService)
	calendarService := service.NewCalendarService(calendarRepo)

	userHandler := handler.NewUserHandler(userService)
	marketHandler := handler.NewMarketHandler(marketService)
	assetHandler := handler.NewAssetHandler(assetService)
	calendarHandler := handler.NewCalendarHandler(calendarService)

	api := app.Group("/api")

	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Sistem hazir ve iliskisel!"})
	})

	auth := api.Group("/v1/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	market := api.Group("/market")
	market.Get("/rates", marketHandler.GetRates)

	protected := api.Group("/v1", middleware.Protected(cfg.JWTSecret))

	u := protected.Group("/user")
	u.Put("/update", userHandler.Update)
	u.Delete("/delete", userHandler.Delete)

	assets := protected.Group("/assets")
	assets.Post("/update", assetHandler.UpdateBalance)
	assets.Get("/summary", assetHandler.GetSummary)
	assets.Get("/transactions", assetHandler.GetTransactions)
	assets.Get("/receipt/full", assetHandler.GetFullPortfolioReceipt)
	assets.Get("/receipt/:id", assetHandler.GetReceipt)
	assets.Get("/export/excel", assetHandler.GetExcel)

	calendar := protected.Group("/calendar")
	calendar.Post("/remind", calendarHandler.AddEvent)
	calendar.Get("/list", calendarHandler.ListReminders)
	calendar.Delete("/:id", calendarHandler.DeleteEvent)
}
