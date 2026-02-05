package main

import (
	"log"

	"github.com/AhmedSelimYildirim/fin-track-pro/internal/handlers"
	"github.com/AhmedSelimYildirim/fin-track-pro/internal/repository" // Burayı ekledik
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Veritabanını başlat
	repository.ConnectDB()

	app := fiber.New()
	handlers.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
