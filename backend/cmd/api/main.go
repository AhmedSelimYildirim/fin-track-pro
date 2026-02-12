package main

import (
	"context"
	_ "fin-track-pro/docs"
	"fin-track-pro/internal/infrastructure/config"
	"fin-track-pro/internal/infrastructure/database"
	"fin-track-pro/internal/infrastructure/redis"
	"fin-track-pro/internal/model"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/router"
	"fin-track-pro/internal/service"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

func main() {
	database.ConnectDB()
	redis.ConnectRedis()
	ctx := context.Background()

	// MarketHistory tablosunu da ekliyoruz
	modelToCreate := []interface{}{
		(*model.User)(nil),
		(*model.Asset)(nil),
		(*model.Transaction)(nil),
		(*model.Reminder)(nil),
		(*model.MarketHistory)(nil),
	}

	for _, m := range modelToCreate {
		_, err := database.DB.NewCreateTable().
			Model(m).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			log.Fatalf("Tablo olusturma hatasi: %v", err)
		}
	}
	fmt.Println("✅ Veritabani tablolari kontrol edildi.")

	cfg := config.LoadConfig()
	marketRepo := repository.NewMarketRepository(database.DB)
	marketService := service.NewMarketService(cfg, redis.Client, marketRepo)

	c := cron.New()
	_, err := c.AddFunc("0 17 * * *", func() {
		log.Println("⏳ Gunluk piyasa verileri kaydediliyor...")
		if err := marketService.SaveDailyRates(); err != nil {
			log.Printf("❌ Kur kaydetme hatasi: %v", err)
		} else {
			log.Println("✅ Gunluk kurlar veritabanina kaydedildi.")
		}
	})
	if err != nil {
		log.Fatal("Cron başlatılamadı:", err)
	}
	c.Start()
	fmt.Println("🕒 Cron servisi aktif (Her gun 17:00)")

	app := fiber.New(fiber.Config{
		AppName:      "FinTrack Pro v2.0",
		ServerHeader: "Fiber",
	})

	router.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Sunucu kapatiliyor...")
	c.Stop()
	_ = app.Shutdown()
	log.Println("FinTrack Pro durduruldu.")
}
