package service

import (
	"bytes"
	"errors"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/repository"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type AssetService struct {
	repo          *repository.AssetRepository
	marketService *MarketService
}

func NewAssetService(repo *repository.AssetRepository, market *MarketService) *AssetService {
	return &AssetService{repo: repo, marketService: market}
}

func (s *AssetService) AddAsset(asset *models.Asset) error {
	return s.repo.Create(asset)
}

func (s *AssetService) UpdateAssetAmount(userID uint, assetID int64, newAmount float64) error {
	asset, err := s.repo.GetByID(assetID)
	if err != nil {
		return err
	}
	if uint(asset.UserID) != userID {
		return errors.New("yetkisiz erisim")
	}
	asset.Amount = newAmount
	return s.repo.Update(asset)
}

func (s *AssetService) RemoveAsset(userID uint, assetID int64) error {
	asset, err := s.repo.GetByID(assetID)
	if err != nil {
		return err
	}
	if uint(asset.UserID) != userID {
		return errors.New("yetkisiz erisim")
	}
	return s.repo.Delete(assetID)
}

func (s *AssetService) GetUserAssets(userID uint) ([]models.Asset, error) {
	return s.repo.GetByUserID(userID)
}

func (s *AssetService) GetAssetByID(userID uint, assetID int64) (*models.Asset, error) {
	asset, err := s.repo.GetByID(assetID)
	if err != nil {
		return nil, err
	}
	if uint(asset.UserID) != userID {
		return nil, errors.New("yetkisiz erisim")
	}
	return asset, nil
}

func (s *AssetService) GetPortfolioSummary(userID uint) (*dto.PortfolioResponse, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	rates, err := s.marketService.GetCurrencyRates()
	if err != nil {
		return nil, errors.New("kurlar alinamadi")
	}

	goldPrice, _ := s.marketService.GetMetalPrice("GOLD")
	silverPrice, _ := s.marketService.GetMetalPrice("SILVER")
	btcPrice, _ := s.marketService.GetCryptoPrice("bitcoin")
	ethPrice, _ := s.marketService.GetCryptoPrice("ethereum")

	groupedAssets := make(map[string]*dto.AssetDetail)
	var total float64

	for _, a := range assets {
		var currentPrice float64
		switch a.Type {
		case "USD":
			currentPrice = rates["USD"]
		case "EUR":
			currentPrice = rates["EUR"]
		case "GOLD":
			currentPrice = goldPrice
		case "SILVER":
			currentPrice = silverPrice
		case "BTC":
			currentPrice = btcPrice
		case "ETH":
			currentPrice = ethPrice
		}

		valueInTL := a.Amount * currentPrice
		total += valueInTL

		if item, exists := groupedAssets[a.Type]; exists {
			item.Amount += a.Amount
			item.ValueInTL += valueInTL
		} else {
			groupedAssets[a.Type] = &dto.AssetDetail{
				Type:         a.Type,
				Amount:       a.Amount,
				CurrentPrice: currentPrice,
				ValueInTL:    valueInTL,
			}
		}
	}

	var details []dto.AssetDetail
	for _, detail := range groupedAssets {
		details = append(details, *detail)
	}

	return &dto.PortfolioResponse{Assets: details, TotalValue: total}, nil
}

func (s *AssetService) ExportToExcel(userID uint) ([]byte, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	f := excelize.NewFile()
	sheet := "Portfoy"
	f.SetSheetName("Sheet1", sheet)
	f.SetCellValue(sheet, "A1", "Varlik Tipi")
	f.SetCellValue(sheet, "B1", "Miktar")
	f.SetCellValue(sheet, "C1", "Alim Tarihi")
	f.SetCellValue(sheet, "D1", "Not")
	for i, a := range assets {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), a.Type)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), a.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), a.PurchaseDate.Format("02.01.2006"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), a.Note)
	}
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *AssetService) GenerateReceipt(asset *models.Asset) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "FINTRACK PRO - ISLEM DEKONTU")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Islem ID: %d", asset.ID))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Varlik Tipi: %s", asset.Type))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Miktar: %.2f", asset.Amount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Tarih: %s", asset.PurchaseDate.Format("02.01.2006")))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Not: %s", asset.Note))
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Bu belge Ahmed Selim YILDIRIM tarafindan uretilmistir.")
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
