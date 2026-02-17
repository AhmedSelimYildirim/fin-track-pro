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

func (r *AssetRepository) GetAsset(userID int64, assetType string, variant string) (*model.Asset, error) {
	asset := new(model.Asset)
	err := r.db.NewSelect().
		Model(asset).
		Where("user_id = ? AND type = ? AND variant = ?", userID, assetType, variant).
		Scan(context.Background())
	return asset, err
}

func (r *AssetRepository) GetByUserID(userID int64) ([]model.Asset, error) {
	var assets []model.Asset
	err := r.db.NewSelect().
		Model(&assets).
		Where("user_id = ?", userID).
		Order("type ASC").
		Order("variant ASC").
		Scan(context.Background())
	return assets, err
}

func (r *AssetRepository) GetTransactionsByUserID(userID int64) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.db.NewSelect().
		Model(&txs).
		Relation("Asset").
		Where("asset.user_id = ?", userID).
		Order("transaction_date DESC").
		Scan(context.Background())
	return txs, err
}

func (r *AssetRepository) GetTransactionByID(txID int64, userID int64) (*model.Transaction, error) {
	tx := new(model.Transaction)
	err := r.db.NewSelect().
		Model(tx).
		Relation("Asset").
		Where("t.id = ?", txID).
		Where("asset.user_id = ?", userID).
		Scan(context.Background())
	return tx, err
}

func (r *AssetRepository) UpdateWithLog(asset *model.Asset, tx *model.Transaction) error {
	ctx := context.Background()
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, bunTx bun.Tx) error {
		_, err := bunTx.NewInsert().
			Model(asset).
			On("CONFLICT (user_id, type, variant) DO UPDATE").
			Set("amount = EXCLUDED.amount").
			Set("updated_at = current_timestamp").
			Returning("id").
			Exec(ctx)
		if err != nil {
			return err
		}

		tx.AssetID = asset.ID
		tx.UserID = asset.UserID
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
