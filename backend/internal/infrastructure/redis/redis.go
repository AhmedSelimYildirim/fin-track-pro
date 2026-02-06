package redis

import (
	"context"
	"fin-track-pro/internal/infrastructure/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func ConnectRedis() {
	cfg := config.LoadConfig()

	if cfg.RedisURL != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			Client = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		} else {
			Client = redis.NewClient(opt)
		}
	} else {
		Client = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}

	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis baglantisi kurulamadi:", err)
	} else {
		fmt.Println("Redis baglantisi aktif 🟢")

		err := Client.FlushAll(ctx).Err()
		if err != nil {
			fmt.Println("Redis sifirlanirken hata olustu:", err)
		} else {
			fmt.Println("Redis verileri basariyla SIFIRLANDI! 🗑️")
		}
	}
}
