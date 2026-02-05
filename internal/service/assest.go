package service

import (
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
)

type AssetService struct {
	repo *repository.AssetRepository
}

func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) AddAsset(asset *models.Asset) error {
	return s.repo.Create(asset)
}

func (s *AssetService) GetUserAssets(userID uint) ([]models.Asset, error) {
	return s.repo.GetByUserID(userID)
}
