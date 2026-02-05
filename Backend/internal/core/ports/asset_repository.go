package ports

import (
	"fin-track-pro/internal/core/models"
)

type AssetRepository interface {
	Create(asset *models.Asset) error
	GetByUserID(userID uint) ([]models.Asset, error)
}
