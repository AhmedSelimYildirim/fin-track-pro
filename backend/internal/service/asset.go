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
	ayar := 0
	if req.Type == "GOLD" {
		ayar = 24
		if req.Ayar > 0 {
			ayar = req.Ayar
		}
	}

	asset, err := s.repo.GetAsset(int64(userID), req.Type, ayar)
	if err != nil {
		asset = &model.Asset{UserID: int64(userID), Type: req.Type, Amount: 0, Ayar: ayar}
	}

	rawPrice, err := s.getCurrentPriceInTRY(req.Type)
	if err != nil || rawPrice <= 0 {
		return errors.New("piyasa verisi su an alinamiyor, lutfen sonra tekrar deneyin")
	}

	unitPrice := rawPrice
	if req.Type == "GOLD" && ayar < 24 {
		unitPrice = rawPrice * (float64(ayar) / 24.0)
	}

	if req.Action == "add" {
		asset.Amount += req.Amount
		asset.TotalCost += req.Amount * unitPrice
	} else if req.Action == "subtract" {
		if asset.Amount < req.Amount {
			return errors.New("yetersiz bakiye")
		}
		if asset.Amount > 0 {
			asset.TotalCost -= req.Amount * (asset.TotalCost / asset.Amount)
		}
		asset.Amount -= req.Amount
	}

	tx := &model.Transaction{
		UserID:          int64(userID),
		Type:            req.Action,
		AssetType:       req.Type,
		Amount:          req.Amount,
		Price:           unitPrice,
		Ayar:            ayar,
		TransactionDate: time.Now(),
	}

	return s.repo.UpdateWithLog(asset, tx)
}

func (s *AssetService) getCurrentPriceInTRY(assetType string) (float64, error) {
	rates, err := s.marketService.GetCurrencyRates()
	if err != nil && assetType != "GOLD" && assetType != "SILVER" {
		return 0, err
	}

	switch assetType {
	case "USD":
		return rates["USD"], nil
	case "EUR":
		return rates["EUR"], nil
	case "GOLD":
		return s.marketService.GetMetalPrice("GOLD")
	case "SILVER":
		return s.marketService.GetMetalPrice("SILVER")
	case "TRY":
		return 1.0, nil
	default:
		return 1.0, nil
	}
}

func (s *AssetService) GetPortfolioSummary(userID uint, baseCurrency string, targetAyar int) (*dto.PortfolioResponse, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	basePriceInTRY, _ := s.getCurrentPriceInTRY(baseCurrency)
	if baseCurrency == "GOLD" && targetAyar > 0 {
		basePriceInTRY *= (float64(targetAyar) / 24.0)
	}
	if basePriceInTRY <= 0 {
		basePriceInTRY = 1
	}

	var details []dto.AssetResponse
	var totalValue float64

	type tempItem struct {
		Type    string
		Amount  float64
		Ayar    int
		ValBase float64
		Rate    float64
	}
	var items []tempItem

	for _, a := range assets {
		rawPriceTRY, err := s.getCurrentPriceInTRY(a.Type)
		if err != nil || rawPriceTRY <= 0 {
			continue
		}
		if a.Type == "GOLD" && a.Ayar > 0 {
			rawPriceTRY *= (float64(a.Ayar) / 24.0)
		}

		rate := rawPriceTRY / basePriceInTRY
		val := a.Amount * rate
		totalValue += val
		items = append(items, tempItem{Type: a.Type, Amount: a.Amount, Ayar: a.Ayar, ValBase: val, Rate: rate})
	}

	for _, it := range items {
		alloc := 0.0
		if totalValue > 0 {
			alloc = (it.ValBase / totalValue) * 100
		}
		details = append(details, dto.AssetResponse{
			Type:         it.Type,
			Amount:       it.Amount,
			Ayar:         it.Ayar,
			CurrentPrice: it.Rate,
			ValueInBase:  it.ValBase,
			Allocation:   alloc,
		})
	}

	return &dto.PortfolioResponse{
		Assets:     details,
		TotalValue: totalValue,
		BaseAsset:  baseCurrency,
	}, nil
}

func (s *AssetService) GenerateTransactionReceipt(tx *model.Transaction, baseCurrency string, targetAyar int) ([]byte, error) {
	basePriceInTRY, _ := s.getCurrentPriceInTRY(baseCurrency)
	if baseCurrency == "GOLD" && targetAyar > 0 {
		basePriceInTRY *= (float64(targetAyar) / 24.0)
	}
	if basePriceInTRY <= 0 {
		basePriceInTRY = 1
	}

	convertedPrice := tx.Price / basePriceInTRY
	userName, _ := s.repo.GetUserName(tx.UserID)
	txTypeTr := "Ekleme"
	if tx.Type == "subtract" {
		txTypeTr = "Cikarma"
	}

	priceFormat := "%.2f"
	if convertedPrice < 0.0001 {
		priceFormat = "%.8f"
	}

	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "FINTRACK PRO - ISLEM DEKONTU")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Varlik: %s", tx.AssetType))
	pdf.Ln(8)
	if tx.AssetType == "GOLD" && tx.Ayar > 0 {
		pdf.Cell(0, 10, fmt.Sprintf("Ayar: %d Ayar", tx.Ayar))
		pdf.Ln(8)
	}
	pdf.Cell(0, 10, fmt.Sprintf("Islem Tipi: %s", txTypeTr))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Miktar: %.4f", tx.Amount))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Birim Fiyat: "+priceFormat+" %s", convertedPrice, baseCurrency))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Toplam Tutar: %.4f %s", tx.Amount*convertedPrice, baseCurrency))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Islem Tarihi: %s", tx.TransactionDate.Format("02.01.2006")))
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Bu belge %s tarafindan uretilmistir.", userName))

	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateFullPortfolioReceipt(userID uint, baseCurrency string, targetAyar int) ([]byte, error) {
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, targetAyar)
	if err != nil {
		return nil, err
	}
	userName, _ := s.repo.GetUserName(int64(userID))

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 15, "FINTRACK PRO - GENEL PORTFOY DEKONTU")
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Varlik")
	pdf.Cell(40, 10, "Miktar")
	pdf.Cell(50, 10, "Birim Deger")
	pdf.Cell(50, 10, fmt.Sprintf("Toplam (%s)", baseCurrency))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for _, a := range summary.Assets {
		pFormat := "%.2f"
		if a.CurrentPrice < 0.0001 {
			pFormat = "%.8f"
		}
		vName := a.Type
		if a.Type == "GOLD" && a.Ayar > 0 {
			vName = fmt.Sprintf("GOLD (%dK)", a.Ayar)
		}
		pdf.Cell(40, 8, vName)
		pdf.Cell(40, 8, fmt.Sprintf("%.4f", a.Amount))
		pdf.Cell(50, 8, fmt.Sprintf(pFormat, a.CurrentPrice))
		pdf.Cell(50, 8, fmt.Sprintf("%.4f", a.ValueInBase))
		pdf.Ln(8)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("TOPLAM PORTFOY DEGERI: %.4f %s", summary.TotalValue, baseCurrency))
	pdf.Ln(15)

	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Rapor Tarihi: %s | Kullanici: %s", time.Now().Format("02.01.2006 15:04"), userName))

	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GetUserTransactions(userID uint) ([]model.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

func (s *AssetService) GetTransactionByID(userID uint, txID int64) (*model.Transaction, error) {
	return s.repo.GetTransactionByID(txID, int64(userID))
}
