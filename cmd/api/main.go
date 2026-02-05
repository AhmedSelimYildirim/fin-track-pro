package main

import (
	"log"

	"fin-track-pro/internal/handlers"
	"fin-track-pro/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	repository.ConnectDB()
	repository.ConnectRedis()

	app := fiber.New(fiber.Config{
		AppName: "FinTrack Pro v1.0",
	})

	handlers.SetupRoutes(app)

	log.Println("🚀 Sunucu http://localhost:3000 adresinde Ahmed Selim için hazır!")
	log.Fatal(app.Listen(":3000"))
}
