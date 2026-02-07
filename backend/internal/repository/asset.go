package repository

import (
	"context"
	"fin-track-pro/internal/model"

	"github.com/uptrace/bun"
)

type AssetRepository struct {
	db *bun.DB
}

func NewAssetRepository(db *bun.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(asset *model.Asset) error {
	_, err := r.db.NewInsert().Model(asset).Exec(context.Background())
	return err
}

func (r *AssetRepository) GetAsset(userID int64, assetType string, ayar int) (*model.Asset, error) {
	asset := new(model.Asset)
	err := r.db.NewSelect().
		Model(asset).
		Where("user_id = ? AND type = ? AND ayar = ?", userID, assetType, ayar).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (r *AssetRepository) GetByUserID(userID uint) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.NewSelect().
		Model(&assets).
		Where("user_id = ?", userID).
		Order("type ASC").
		Scan(context.Background())
	return assets, err
}

func (r *AssetRepository) GetTransactionsByUserID(userID uint) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.NewSelect().
		Model(&txs).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(context.Background())
	return txs, err
}

func (r *AssetRepository) GetTransactionByID(txID int64, userID int64) (*model.Transaction, error) {
	tx := new(model.Transaction)
	err := r.db.NewSelect().
		Model(tx).
		Where("id = ? AND user_id = ?", txID, userID).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *AssetRepository) UpdateWithLog(asset *model.Asset, tx *model.Transaction) error {
	ctx := context.Background()
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, bunTx bun.Tx) error {
		_, err := bunTx.NewInsert().
			Model(asset).
			On("CONFLICT (user_id, type, ayar) DO UPDATE").
			Set("amount = EXCLUDED.amount").
			Set("total_cost = EXCLUDED.total_cost").
			Exec(ctx)
		if err != nil {
			return err
		}
		if _, err := bunTx.NewInsert().Model(tx).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}

func (r *AssetRepository) GetUserName(userID int64) (string, error) {
	var user model.User
	err := r.db.NewSelect().
		Model(&user).
		Column("username").
		Where("id = ?", userID).
		Scan(context.Background())
	if err != nil {
		return "Kullanici", nil
	}
	return user.Username, nil
}
