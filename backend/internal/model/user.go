package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64     `bun:",pk,autoincrement" json:"id"`
	FullName      string    `bun:"full_name,notnull" json:"full_name"`
	Username      string    `bun:"username,unique,notnull" json:"username"`
	Email         string    `bun:"email,unique,notnull" json:"email"`
	Password      string    `bun:"password,notnull" json:"-"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt     time.Time `bun:",soft_delete,nullzero" json:"-"`
}
