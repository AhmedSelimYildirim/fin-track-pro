package main

import (
	"context"
	_ "fin-track-pro/docs"
	"fin-track-pro/internal/infrastructure/database"
	"fin-track-pro/internal/infrastructure/redis"
	"fin-track-pro/internal/model"
	"fin-track-pro/internal/router"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()
	redis.ConnectRedis()
	ctx := context.Background()

	modelToCreate := []interface{}{
		(*model.User)(nil),
		(*model.Asset)(nil),
		(*model.Transaction)(nil),
		(*model.Reminder)(nil),
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
	fmt.Println("✅ Veritabani kontrol edildi, eksik tablolar olusturuldu.")

	app := fiber.New(fiber.Config{
		AppName:      "FinTrack Pro v1.0",
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Sunucu kapatiliyor...")
	_ = app.Shutdown()
	log.Println("FinTrack Pro durduruldu. Gorusuruz !")
}
