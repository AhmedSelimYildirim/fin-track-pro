package dto

import "time"

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AssetCreateRequest struct {
	Type         string     `json:"type" validate:"required"`
	Amount       float64    `json:"amount" validate:"required,gt=0"`
	Cost         float64    `json:"cost"`
	Note         string     `json:"note"`
	PurchaseDate time.Time  `json:"purchase_date"`
	CalendarNote *string    `json:"calendar_note"`
	ReminderDate *time.Time `json:"reminder_date"`
}

type AssetDetail struct {
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	CurrentPrice float64 `json:"current_price"`
	ValueInTL    float64 `json:"value_in_tl"`
}

type AssetResponse struct {
	ID           int64   `json:"id"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	CurrentPrice float64 `json:"current_price"`
	ValueInTL    float64 `json:"value_in_tl"`
	Message      string  `json:"message"`
}

type PortfolioResponse struct {
	Assets     []AssetDetail `json:"assets"`
	TotalValue float64       `json:"total_value"`
}
