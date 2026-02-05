package repository

import (
	"fin-track-pro/internal/core/models"

	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(asset *models.Asset) error {
	return r.db.Create(asset).Error
}

func (r *AssetRepository) GetByUserID(userID uint) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.Where("user_id = ?", userID).Find(&assets).Error
	return assets, err
}
