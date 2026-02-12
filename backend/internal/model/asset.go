package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"`

	ID        int64     `bun:",pk,autoincrement" json:"id"`
	UserID    int64     `bun:"user_id,notnull" json:"user_id"`
	Type      string    `bun:"type,notnull" json:"type"`
	Variant   string    `bun:"variant,notnull" json:"variant"`
	Amount    float64   `bun:"amount,notnull,default:0" json:"amount"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	User         *User          `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	Transactions []*Transaction `bun:"rel:has-many,join:id=asset_id" json:"transactions,omitempty"`
}
