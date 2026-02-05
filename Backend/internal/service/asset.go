package service

import (
	"bytes"
	"errors"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/repository"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type AssetService struct {
	repo          *repository.AssetRepository
	marketService *MarketService
}

func NewAssetService(repo *repository.AssetRepository, market *MarketService) *AssetService {
	return &AssetService{repo: repo, marketService: market}
}

func (s *AssetService) ManageBalance(userID uint, assetType string, amount float64, action string) error {
	asset, err := s.repo.GetByType(int64(userID), assetType)
	if err != nil {
		asset = &models.Asset{UserID: int64(userID), Type: assetType, Amount: 0}
		if err := s.repo.Create(asset); err != nil {
			return err
		}
	}

	currentPrice := s.getCurrentPrice(assetType)

	if action == "add" {
		asset.Amount += amount
	} else if action == "subtract" {
		if asset.Amount < amount {
			return errors.New("yetersiz bakiye ")
		}
		asset.Amount -= amount
	} else {
		return errors.New("gecersiz islem")
	}

	tx := &models.Transaction{
		UserID:    int64(userID),
		Type:      action,
		AssetType: assetType,
		Amount:    amount,
		Price:     currentPrice,
	}

	return s.repo.UpdateWithLog(asset, tx)
}

func (s *AssetService) getCurrentPrice(assetType string) float64 {
	rates, _ := s.marketService.GetCurrencyRates()
	switch assetType {
	case "USD":
		return rates["USD"]
	case "GOLD":
		p, _ := s.marketService.GetMetalPrice("GOLD")
		return p
	case "SILVER":
		p, _ := s.marketService.GetMetalPrice("SILVER")
		return p
	case "BTC":
		p, _ := s.marketService.GetCryptoPrice("bitcoin")
		return p
	default:
		return 1.0
	}
}

func (s *AssetService) GetPortfolioSummary(userID uint) (*dto.PortfolioResponse, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	var details []dto.AssetDetail
	var total float64
	for _, a := range assets {
		price := s.getCurrentPrice(a.Type)
		val := a.Amount * price
		total += val
		details = append(details, dto.AssetDetail{
			Type: a.Type, Amount: a.Amount, CurrentPrice: price, ValueInTL: val,
		})
	}
	return &dto.PortfolioResponse{Assets: details, TotalValue: total}, nil
}

func (s *AssetService) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

func (s *AssetService) GetTransactionByID(userID uint, txID int64) (*models.Transaction, error) {
	return s.repo.GetTransactionByID(txID, int64(userID))
}

func (s *AssetService) GenerateTransactionReceipt(tx *models.Transaction) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "FINTRACK PRO - ISLEM DEKONTU")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Islem No: %d", tx.ID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Varlik: %s", tx.AssetType))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Islem Tipi: %s", tx.Type))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Miktar: %.4f", tx.Amount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Birim Fiyat: %.2f TL", tx.Price))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Tarih: %s", tx.CreatedAt.Format("02.01.2006 15:04")))
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Bu belge  YILDIRIM tarafindan uretilmistir.")
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
