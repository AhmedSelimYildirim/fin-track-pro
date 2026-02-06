package main

import (
	"context"
	_ "fin-track-pro/docs"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/router"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	repository.ConnectDB()
	repository.ConnectRedis()
	ctx := context.Background()

	modelsToCreate := []interface{}{
		(*models.User)(nil),
		(*models.Asset)(nil),
		(*models.Transaction)(nil),
		(*models.Reminder)(nil),
	}

	for _, model := range modelsToCreate {
		_, err := repository.DB.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			log.Fatalf("Tablo olusturma hatasi: %v", err)
		}
	}

	app := fiber.New(fiber.Config{
		AppName:      "FinTrack Pro v1.0",
		ServerHeader: "Fiber",
	})

	router.SetupRoutes(app)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		if err := app.Listen(":" + port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Sunucu kapatılıyor...")
	_ = app.Shutdown()
	log.Println("FinTrack Pro durduruldu. Görüşürüz !")
}
