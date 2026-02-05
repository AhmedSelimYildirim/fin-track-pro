package main

import (
	"context"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/handlers"
	repository2 "fin-track-pro/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "fin-track-pro/docs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	repository2.ConnectDB()
	repository2.ConnectRedis()

	ctx := context.Background()

	modelsToCreate := []interface{}{
		(*models.User)(nil),
		(*models.Asset)(nil),
		(*models.Calendar)(nil),
	}

	for _, model := range modelsToCreate {
		_, err := repository2.DB.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			log.Fatalf("Tablo olusturma hatasi: %v", err)
		}
	}
	log.Println("Veritabanı şeması Ahmed Selim için hazırlandı.")

	app := fiber.New(fiber.Config{
		AppName:      "FinTrack Pro v1.0",
		ServerHeader: "Fiber",
	})

	handlers.SetupRoutes(app)

	go func() {
		log.Println("Hatırlatıcı Worker servisi başlatıldı...")
		for {
			time.Sleep(1 * time.Minute)
		}
	}()

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Sunucu kapatılıyor...")

	_ = app.Shutdown()
	log.Println("FinTrack Pro durduruldu. Görüşürüz Ahmed Selim!")
}
