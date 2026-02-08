package redis

import (
	"context"
	"fin-track-pro/internal/infrastructure/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func ConnectRedis() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	if cfg.RedisURL != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			log.Printf("Redis URL parse hatasi: %v", err)
			Client = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		} else {
			Client = redis.NewClient(opt)
		}
	} else {
		Client = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis baglantisi kurulamadi 🔴:", err)
	} else {
		fmt.Println("Redis baglantisi aktif 🟢")

		err := Client.FlushAll(ctx).Err()
		if err != nil {
			fmt.Println("Redis temizleme hatasi:", err)
		} else {
			fmt.Println("🚀 Redis hafizasi tamamen temizlendi!")
		}
	}
}
