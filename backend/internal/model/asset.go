package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"`
	ID            int64     `bun:",pk,autoincrement" json:"id"`
	UserID        int64     `bun:"user_id,notnull,unique:user_asset_ayar" json:"user_id"`
	Type          string    `bun:"type,notnull,unique:user_asset_ayar" json:"type"`
	Amount        float64   `bun:"amount,notnull,default:0" json:"amount"`
	Ayar          int       `bun:"ayar,default:24,unique:user_asset_ayar" json:"ayar"`
	TotalCost     float64   `bun:"total_cost,default:0" json:"total_cost"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
