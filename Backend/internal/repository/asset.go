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

func (r *AssetRepository) GetByType(userID int64, assetType string) (*models.Asset, error) {
	asset := new(models.Asset)
	err := r.db.NewSelect().
		Model(asset).
		Where("user_id = ? AND type = ?", userID, assetType).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (r *AssetRepository) GetByUserID(userID uint) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.NewSelect().
		Model(&assets).
		Where("user_id = ?", userID).
		Scan(context.Background())
	return assets, err
}

func (r *AssetRepository) GetTransactionByID(txID int64, userID int64) (*models.Transaction, error) {
	tx := new(models.Transaction)
	err := r.db.NewSelect().
		Model(tx).
		Where("id = ? AND user_id = ?", txID, userID).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *AssetRepository) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := r.db.NewSelect().
		Model(&txs).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(context.Background())
	return txs, err
}

func (r *AssetRepository) UpdateWithLog(asset *models.Asset, tx *models.Transaction) error {
	ctx := context.Background()
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, bunTx bun.Tx) error {
		if _, err := bunTx.NewUpdate().Model(asset).WherePK().Exec(ctx); err != nil {
			return err
		}
		if _, err := bunTx.NewInsert().Model(tx).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}
