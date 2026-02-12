package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Reminder struct {
	bun.BaseModel `bun:"table:reminders,alias:r"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID     int64     `bun:"user_id,notnull" json:"user_id"`
	Title      string    `bun:"title,notnull" json:"title"`
	TargetDate time.Time `bun:"target_date,notnull" json:"target_date"`
	CreatedAt  time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`

	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}
