package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:",pk,autoincrement" json:"id"`
	Username      string    `bun:"username,unique,notnull" json:"username"`
	Email         string    `bun:"email,unique,notnull" json:"email"`
	Password      string    `bun:"password,notnull" json:"-"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero" json:"-"`
}

type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"`
	ID            int64     `bun:",pk,autoincrement" json:"id"`
	UserID        int64     `bun:"user_id,notnull" json:"user_id"`
	Type          string    `bun:"type,unique:user_asset_type,notnull" json:"type"`
	Amount        float64   `bun:"amount,notnull,default:0" json:"amount"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`
	ID            int64     `bun:",pk,autoincrement" json:"id"`
	UserID        int64     `bun:"user_id,notnull" json:"user_id"`
	Type          string    `bun:"type,notnull" json:"type"`
	AssetType     string    `bun:"asset_type,notnull" json:"asset_type"`
	Amount        float64   `bun:"amount,notnull" json:"amount"`
	Price         float64   `bun:"price" json:"price"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type Reminder struct {
	bun.BaseModel `bun:"table:reminders,alias:r"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID        int64     `bun:"user_id" json:"user_id"`
	Title         string    `bun:"title" json:"title"`
	TargetDate    time.Time `bun:"target_date" json:"target_date"`
	IsSent        bool      `bun:"is_sent,default:false" json:"is_sent"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
}
