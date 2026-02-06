package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel   `bun:"table:transactions,alias:t"`
	ID              int64     `bun:",pk,autoincrement" json:"id"`
	UserID          int64     `bun:"user_id,notnull" json:"user_id"`
	Type            string    `bun:"type,notnull" json:"type"`
	AssetType       string    `bun:"asset_type,notnull" json:"asset_type"`
	Amount          float64   `bun:"amount,notnull" json:"amount"`
	Price           float64   `bun:"price" json:"price"`
	Ayar            int       `bun:"ayar,default:24" json:"ayar"`
	TransactionDate time.Time `bun:"transaction_date,notnull" json:"transaction_date"`
	CreatedAt       time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}
