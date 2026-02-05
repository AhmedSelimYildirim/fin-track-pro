package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type MarketRepository struct {
	redis *redis.Client
}

func NewMarketRepository(client *redis.Client) *MarketRepository {
	return &MarketRepository{redis: client}
}

func (r *MarketRepository) SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.redis.Set(ctx, key, value, expiration).Err()
}

func (r *MarketRepository) GetCache(key string) (string, error) {
	ctx := context.Background()
	return r.redis.Get(ctx, key).Result()
}
