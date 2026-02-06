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

	_, err := Client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis baglantisi kurulamadi")
	} else {
		fmt.Println("Redis baglantisi aktif")
	}
}
