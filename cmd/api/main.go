package main

import (
	"log"

	"github.com/AhmedSelimYildirim/fin-track-pro/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		AppName: "FinTrack Pro v1.0",
	})

	handlers.SetupRoutes(app)

	log.Println("Sunucu 3000 portunda başlatılıyor...")
	log.Fatal(app.Listen(":3000"))
}
