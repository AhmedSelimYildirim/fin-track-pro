package dto

import "time"

type TransactionResponse struct {
	ID              int64     `json:"id"`
	Type            string    `json:"type"`
	AssetType       string    `json:"asset_type"`
	Amount          float64   `json:"amount"`
	Price           float64   `json:"price"`
	ayar            int       `json:"ayar"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
}
