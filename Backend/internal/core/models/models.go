package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Username  string    `bun:"username,unique,notnull" json:"username"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	Password  string    `bun:"password,notnull" json:"-"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"`

	ID           int64      `bun:",pk,autoincrement" json:"id"`
	UserID       int64      `bun:"user_id,notnull" json:"user_id"`
	Type         string     `bun:"type,notnull" json:"type"`
	Amount       float64    `bun:"amount,notnull" json:"amount"`
	Cost         float64    `bun:"cost,notnull" json:"cost"`
	Note         string     `bun:"note" json:"note"`
	PurchaseDate time.Time  `bun:"purchase_date" json:"purchase_date"`
	CalendarNote string     `bun:"calendar_note" json:"calendar_note"`
	ReminderDate *time.Time `bun:"reminder_date" json:"reminder_date"`
	CreatedAt    time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt    time.Time  `bun:",soft_delete,nullzero" json:"-"`
}

type Reminder struct {
	bun.BaseModel `bun:"table:reminders,alias:r"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID     int64     `bun:"user_id" json:"user_id"`
	Title      string    `bun:"title" json:"title"`
	TargetDate time.Time `bun:"target_date" json:"target_date"`
	IsSent     bool      `bun:"is_sent,default:false" json:"is_sent"`
	CreatedAt  time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
}

type Calendar struct {
	bun.BaseModel `bun:"table:calendars,alias:c"`

	ID           int64     `json:"id" bun:",pk,autoincrement"`
	UserID       int64     `json:"user_id" bun:",notnull"`
	EventName    string    `json:"event_name" bun:",notnull"`
	ReminderDate time.Time `json:"reminder_date" bun:",notnull"`
	Note         string    `json:"note"`
	CreatedAt    time.Time `json:"created_at" bun:",default:current_timestamp"`
}
