package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Reminder struct {
	bun.BaseModel `bun:"table:reminders,alias:r"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID        int64     `bun:"user_id" json:"user_id"`
	Title         string    `bun:"title" json:"title"`
	TargetDate    time.Time `bun:"target_date" json:"target_date"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
}
