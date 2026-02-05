package repository

import (
	"github.com/uptrace/bun"
)

type MarketRepository struct {
	db *bun.DB
}

func NewMarketRepository(db *bun.DB) *MarketRepository {
	return &MarketRepository{db: db}
}
