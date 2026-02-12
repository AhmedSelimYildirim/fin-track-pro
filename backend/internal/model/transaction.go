package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID              int64     `bun:",pk,autoincrement" json:"id"`
	AssetID         int64     `bun:"asset_id,notnull" json:"asset_id"`
	UserID          int64     `bun:"user_id,notnull" json:"user_id"`
	Type            string    `bun:"type,notnull" json:"type"`
	Amount          float64   `bun:"amount,notnull" json:"amount"`
	Price           float64   `bun:"price,notnull" json:"price"`
	TransactionDate time.Time `bun:"transaction_date,notnull" json:"transaction_date"`
	CreatedAt       time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Asset *Asset `bun:"rel:belongs-to,join:asset_id=id" json:"asset,omitempty"`
}
