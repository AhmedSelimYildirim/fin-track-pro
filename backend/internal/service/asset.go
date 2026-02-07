package service

import (
	"bytes"
	"errors"
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/model"
	"fin-track-pro/internal/repository"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type AssetService struct {
	repo          *repository.AssetRepository
	marketService *MarketService
}

func NewAssetService(repo *repository.AssetRepository, market *MarketService) *AssetService {
	return &AssetService{repo: repo, marketService: market}
}

func (s *AssetService) ManageBalance(userID uint, req dto.AssetCreateRequest) error {
	ayar := 24
	if req.Type == "GOLD" && req.Ayar > 0 {
		ayar = req.Ayar
	}

	asset, err := s.repo.GetAsset(int64(userID), req.Type, ayar)
	if err != nil {
		asset = &model.Asset{UserID: int64(userID), Type: req.Type, Amount: 0, Ayar: ayar}
	}

	unitPrice := s.getCurrentPriceInTRY(req.Type)

	if req.Action == "add" {
		asset.Amount += req.Amount
		asset.TotalCost += req.Amount * unitPrice
	} else if req.Action == "subtract" {
		if asset.Amount < req.Amount {
			return errors.New("yetersiz bakiye")
		}
		if asset.Amount > 0 {
			avgCost := asset.TotalCost / asset.Amount
			asset.TotalCost -= req.Amount * avgCost
		}
		asset.Amount -= req.Amount
	} else {
		return errors.New("gecersiz islem")
	}

	tx := &model.Transaction{
		UserID: int64(userID), Type: req.Action, AssetType: req.Type,
		Amount: req.Amount, Price: unitPrice, Ayar: ayar, TransactionDate: time.Now(),
	}

	return s.repo.UpdateWithLog(asset, tx)
}

func (s *AssetService) getCurrentPriceInTRY(assetType string) float64 {
	rates, _ := s.marketService.GetCurrencyRates()
	usdToTry := rates["USD"]

	switch assetType {
	case "USD":
		return usdToTry
	case "EUR":
		return rates["EUR"] * usdToTry
	case "GOLD":
		p, _ := s.marketService.GetMetalPrice("GOLD")
		return p
	case "SILVER":
		p, _ := s.marketService.GetMetalPrice("SILVER")
		return p
	case "BTC":
		p, _ := s.marketService.GetCryptoPrice("bitcoin")
		return p
	case "TRY":
		return 1.0
	default:
		return 1.0
	}
}

func (s *AssetService) GetPortfolioSummary(userID uint, baseCurrency string, targetAyar int) (*dto.PortfolioResponse, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	baseCurrencyPriceInTRY := s.getCurrentPriceInTRY(baseCurrency)
	if baseCurrency == "GOLD" && targetAyar > 0 && targetAyar < 24 {
		baseCurrencyPriceInTRY = baseCurrencyPriceInTRY * (float64(targetAyar) / 24.0)
	}

	if baseCurrencyPriceInTRY == 0 {
		baseCurrencyPriceInTRY = 1
	}

	var details []dto.AssetResponse
	var totalValue float64
	var totalCostInBase float64

	for _, a := range assets {
		assetRawPriceInTRY := s.getCurrentPriceInTRY(a.Type)

		adjustedPriceInTRY := assetRawPriceInTRY
		if a.Type == "GOLD" {
			adjustedPriceInTRY = assetRawPriceInTRY * (float64(a.Ayar) / 24.0)
		}

		exchangeRate := adjustedPriceInTRY / baseCurrencyPriceInTRY
		valInBase := a.Amount * exchangeRate
		costInBase := a.TotalCost / baseCurrencyPriceInTRY
		profitLoss := valInBase - costInBase

		profitLossRatio := 0.0
		if costInBase > 0 {
			profitLossRatio = (profitLoss / costInBase) * 100
		}

		totalValue += valInBase
		totalCostInBase += costInBase

		displayAyar := 0
		if a.Type == "GOLD" {
			displayAyar = a.Ayar
		}

		details = append(details, dto.AssetResponse{
			Type:            a.Type,
			Amount:          a.Amount,
			Ayar:            displayAyar,
			CurrentPrice:    exchangeRate,
			ValueInBase:     valInBase,
			ProfitLoss:      profitLoss,
			ProfitLossRatio: profitLossRatio,
		})
	}

	baseAssetLabel := baseCurrency
	if baseCurrency == "GOLD" && targetAyar > 0 {
		baseAssetLabel = fmt.Sprintf("GOLD (%d Ayar)", targetAyar)
	}

	return &dto.PortfolioResponse{
		Assets:          details,
		TotalValue:      totalValue,
		TotalCost:       totalCostInBase,
		TotalProfitLoss: totalValue - totalCostInBase,
		BaseAsset:       baseAssetLabel,
	}, nil
}

func (s *AssetService) GetUserTransactions(userID uint) ([]model.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

func (s *AssetService) GetTransactionByID(userID uint, txID int64) (*model.Transaction, error) {
	return s.repo.GetTransactionByID(txID, int64(userID))
}

func (s *AssetService) GenerateTransactionReceipt(tx *model.Transaction, baseCurrency string, targetAyar int) ([]byte, error) {
	baseCurrencyPriceInTRY := s.getCurrentPriceInTRY(baseCurrency)

	baseLabel := baseCurrency
	if baseCurrency == "GOLD" && targetAyar > 0 && targetAyar < 24 {
		baseCurrencyPriceInTRY = baseCurrencyPriceInTRY * (float64(targetAyar) / 24.0)
		baseLabel = fmt.Sprintf("GOLD (%d Ayar)", targetAyar)
	}

	if baseCurrencyPriceInTRY == 0 {
		baseCurrencyPriceInTRY = 1
	}

	convertedPrice := tx.Price / baseCurrencyPriceInTRY

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
	if tx.Type == "GOLD" && tx.Ayar > 0 {
		pdf.Cell(0, 10, fmt.Sprintf("Ayar: %d Ayar", tx.Ayar))
		pdf.Ln(8)
	}
	pdf.Cell(0, 10, fmt.Sprintf("Islem Tipi: %s", tx.Type))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Miktar: %.4f", tx.Amount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Birim Fiyat: %.2f %s", convertedPrice, baseLabel))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Islem Tarihi: %s", tx.TransactionDate.Format("02.01.2006")))
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Bu belge Ahmed Selim YILDIRIM tarafindan uretilmistir.")
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateFullPortfolioReceipt(userID uint, baseCurrency string, targetAyar int) ([]byte, error) {
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, targetAyar)
	if err != nil {
		return nil, err
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 15, "FINTRACK PRO - GENEL PORTFOY DEKONTU")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(30, 10, "Varlik")
	pdf.Cell(30, 10, "Miktar")
	pdf.Cell(30, 10, "Birim Deger")
	pdf.Cell(40, 10, fmt.Sprintf("Toplam (%s)", summary.BaseAsset))
	pdf.Cell(40, 10, fmt.Sprintf("Kar/Zarar (%s)", summary.BaseAsset))
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	for _, a := range summary.Assets {
		pdf.Cell(30, 8, a.Type)
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", a.Amount))
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", a.CurrentPrice))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.ValueInBase))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.ProfitLoss))
		pdf.Ln(8)
	}
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("TOPLAM DEGER: %.2f %s", summary.TotalValue, summary.BaseAsset))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("TOPLAM KAR/ZARAR: %.2f %s", summary.TotalProfitLoss, summary.BaseAsset))
	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Rapor Tarihi: %s", time.Now().Format("02.01.2006 15:04")))
	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}
