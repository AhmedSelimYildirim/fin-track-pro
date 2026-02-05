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
	UserID uint      `json:"user_id"`
	Type   string    `json:"type"`   // GOLD, USD, BTC, SILVER, ETH vb.
	Amount float64   `json:"amount"` // Elindeki miktar (örneğin 10.5 gram)
	Cost   float64   `json:"cost"`   // Alış fiyatı (Kar/Zarar hesabı için)
	Note   string    `json:"note"`   // "Bugün maaşla aldım" gibi notlar
	Date   time.Time `json:"date"`   // Takvim bazlı takip için tarih
}
