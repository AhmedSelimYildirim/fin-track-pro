package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Asset struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	Type      string  `gorm:"not null"`
	Amount    float64 `gorm:"not null"`
	CostPrice float64
	Date      time.Time
}
