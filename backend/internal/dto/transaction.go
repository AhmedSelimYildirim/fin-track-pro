package dto

import "time"

type TransactionResponse struct {
	ID              int64     `json:"id"`
	AssetID         int64     `json:"asset_id"`
	Type            string    `json:"type"`
	AssetType       string    `json:"asset_type"`
	Variant         string    `json:"variant"`
	Amount          float64   `json:"amount"`
	Price           float64   `json:"price"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
}
