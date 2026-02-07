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
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	Ayar            int     `json:"ayar"`
	CurrentPrice    float64 `json:"current_price"`
	ValueInBase     float64 `json:"value_in_base"`
	ProfitLoss      float64 `json:"profit_loss"`
	ProfitLossRatio float64 `json:"profit_loss_ratio"`
}

type PortfolioResponse struct {
	Assets          []AssetResponse `json:"assets"`
	TotalValue      float64         `json:"total_value"`
	TotalCost       float64         `json:"total_cost"`
	TotalProfitLoss float64         `json:"total_profit_loss"`
	BaseAsset       string          `json:"base_asset"`
}
