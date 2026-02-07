package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"`
	ID            int64     `bun:",pk,autoincrement"`
	UserID        int64     `bun:"user_id,notnull,unique:user_asset_ayar"`
	Type          string    `bun:"type,notnull,unique:user_asset_ayar"`
	Ayar          int       `bun:"ayar,default:24,unique:user_asset_ayar"`
	Amount        float64   `bun:"amount,notnull,default:0"`
	TotalCost     float64   `bun:"total_cost,default:0"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
