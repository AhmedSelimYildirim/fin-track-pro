package repository

import (
	"context"
	"fin-track-pro/internal/model"

	"github.com/uptrace/bun"
)

type MarketRepository struct {
	db *bun.DB
}

func NewMarketRepository(db *bun.DB) *MarketRepository {
	return &MarketRepository{db: db}
}

func (r *MarketRepository) SaveRates(history *model.MarketHistory) error {
	_, err := r.db.NewInsert().Model(history).Exec(context.Background())
	return err
}
