package dto

import "time"

type AssetCreateRequest struct {
	Type            string     `json:"type"`
	Variant         string     `json:"variant"`
	Amount          float64    `json:"amount"`
	Action          string     `json:"action"`
	TransactionDate *time.Time `json:"transaction_date"`
}

type AssetResponse struct {
	ID           int64   `json:"id"`
	Type         string  `json:"type"`
	Variant      string  `json:"variant"`
	Amount       float64 `json:"amount"`
	CurrentPrice float64 `json:"current_price"`
	ValueInBase  float64 `json:"value_in_base"`
	Allocation   float64 `json:"allocation"`
}

type PortfolioResponse struct {
	Assets     []AssetResponse `json:"assets"`
	TotalValue float64         `json:"total_value"`
	BaseAsset  string          `json:"base_asset"`
}
