package main

import (
	"context"
	_ "fin-track-pro/docs"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/router"
	"fin-track-pro/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
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

	calendarRepo := repository.NewCalendarRepository(repository.DB)
	calendarService := service.NewCalendarService(calendarRepo)

	app := fiber.New(fiber.Config{
		AppName:      "FinTrack Pro v1.0",
		ServerHeader: "Fiber",
	})

	router.SetupRoutes(app)

	go func() {
		log.Println("Hatırlatıcı Worker servisi aktif edildi...")
		for {
			calendarService.ProcessPendingReminders()
			time.Sleep(1 * time.Minute)
		}
	}()

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
