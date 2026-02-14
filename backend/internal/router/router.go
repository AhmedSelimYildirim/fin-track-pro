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
		AllowOrigins:     "https://fin-track-pro.onrender.com, https://fin-track-pro-1.onrender.com, http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Currency, X-Ayar",
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("FinTrack Pro Backend is Running!")
	})

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
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	v1.Get("/market/rates", marketHandler.GetRates)

	user := v1.Group("/user", middleware.Protected(cfg.JWTSecret))
	user.Put("/update", userHandler.Update)
	user.Delete("/delete", userHandler.Delete)

	assets := v1.Group("/assets", middleware.Protected(cfg.JWTSecret))
	assets.Post("/update", assetHandler.UpdateBalance)
	assets.Get("/summary", assetHandler.GetSummary)
	assets.Get("/transactions", assetHandler.GetTransactions)
	assets.Get("/receipt/full", assetHandler.GetFullPortfolioReceipt)
	assets.Get("/receipt/:id", assetHandler.GetReceipt)
	assets.Get("/export/excel", assetHandler.GetExcel)

	calendar := v1.Group("/calendar", middleware.Protected(cfg.JWTSecret))
	calendar.Post("/remind", calendarHandler.AddEvent)
	calendar.Get("/list", calendarHandler.ListReminders)
	calendar.Delete("/:id", calendarHandler.DeleteEvent)
}
