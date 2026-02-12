package model

import (
	"time"

	"github.com/uptrace/bun"
)

type MarketHistory struct {
	bun.BaseModel `bun:"table:market_history,alias:mh"`

	ID        int64              `bun:",pk,autoincrement" json:"id"`
	Date      time.Time          `bun:"date,notnull" json:"date"`
	Rates     map[string]float64 `bun:"rates,type:jsonb" json:"rates"`
	CreatedAt time.Time          `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}
