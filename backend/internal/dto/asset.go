package dto

import "time"

type AssetCreateRequest struct {
	Type            string     `json:"type" validate:"required"`
	Amount          float64    `json:"amount" validate:"required,gt=0"`
	Action          string     `json:"action" validate:"required,oneof=add subtract"`
	Ayar            int        `json:"ayar"`
	TransactionDate *time.Time `json:"transaction_date"`
}

type AssetResponse struct {
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	Ayar         int     `json:"ayar,omitempty"`
	CurrentPrice float64 `json:"current_price"`
	ValueInBase  float64 `json:"value_in_base"`
	Allocation   float64 `json:"allocation"`
}

type PortfolioResponse struct {
	Assets     []AssetResponse `json:"assets"`
	TotalValue float64         `json:"total_value"`
	BaseAsset  string          `json:"base_asset"`
}
