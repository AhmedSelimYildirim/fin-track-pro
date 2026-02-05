package repository

import (
	"context"
	"fin-track-pro/internal/core/models"

	"github.com/uptrace/bun"
)

type AssetRepository struct {
	db *bun.DB
}

func NewAssetRepository(db *bun.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(asset *models.Asset) error {
	_, err := r.db.NewInsert().Model(asset).Exec(context.Background())
	return err
}

func (r *AssetRepository) GetByUserID(userID uint) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.NewSelect().
		Model(&assets).
		Where("user_id = ?", userID).
		Order("purchase_date DESC").
		Scan(context.Background())
	return assets, err
}

func (r *AssetRepository) GetByID(id int64) (*models.Asset, error) {
	asset := new(models.Asset)
	err := r.db.NewSelect().Model(asset).Where("id = ?", id).Scan(context.Background())
	return asset, err
}

func (r *AssetRepository) Update(asset *models.Asset) error {
	_, err := r.db.NewUpdate().Model(asset).WherePK().Exec(context.Background())
	return err
}

func (r *AssetRepository) Delete(id int64) error {
	_, err := r.db.NewDelete().Model((*models.Asset)(nil)).Where("id = ?", id).Exec(context.Background())
	return err
}
