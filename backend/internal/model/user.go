package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Username  string    `bun:"username,notnull" json:"username"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	Password  string    `bun:"password,notnull" json:"-"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`

	Assets    []*Asset    `bun:"rel:has-many,join:id=user_id" json:"assets,omitempty"`
	Reminders []*Reminder `bun:"rel:has-many,join:id=user_id" json:"reminders,omitempty"`
}
