package main

import (
	"context"
	_ "fin-track-pro/docs"
	"fin-track-pro/internal/infrastructure/database"
	"fin-track-pro/internal/infrastructure/redis"
	"fin-track-pro/internal/model"
	"fin-track-pro/internal/router"
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

	// NOT: Silme (Drop) kodlarını kaldırdık. Artık veriler kalıcı!

	modelToCreate := []interface{}{
		(*model.User)(nil),
		(*model.Asset)(nil),
		(*model.Transaction)(nil),
		(*model.Reminder)(nil),
	}

	for _, m := range modelToCreate {
		_, err := database.DB.NewCreateTable().
			Model(m).
			IfNotExists(). // Sadece tablo yoksa oluşturur, varsa dokunmaz.
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
