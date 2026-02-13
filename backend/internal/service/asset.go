package service

import (
	"bytes"
	"errors"
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/model"
	"fin-track-pro/internal/repository"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

var GoldFactors = map[string]float64{
	"GRAM_24":    1.0,
	"GRAM_22":    0.916,
	"GRAM_18":    0.750,
	"GRAM_14":    0.585,
	"GRAM_8":     0.333,
	"GRAM_4":     0.166,
	"CEYREK":     1.605,
	"YARIM":      3.21,
	"TAM":        6.42,
	"CUMHURIYET": 6.61,
	"GREMSE":     16.05,
}

type AssetService struct {
	repo          *repository.AssetRepository
	marketService *MarketService
}

func NewAssetService(repo *repository.AssetRepository, market *MarketService) *AssetService {
	return &AssetService{repo: repo, marketService: market}
}

func (s *AssetService) tr(text string) string {
	replacer := strings.NewReplacer("ğ", "g", "Ğ", "G", "ü", "u", "Ü", "U", "ş", "s", "Ş", "S", "ı", "i", "İ", "I", "ö", "o", "Ö", "O", "ç", "c", "Ç", "C")
	return replacer.Replace(text)
}

func (s *AssetService) ManageBalance(userID int64, req dto.AssetCreateRequest) error {
	variant := req.Variant
	if variant == "" || variant == "STANDARD" {
		if req.Type == "GOLD" {
			variant = "GRAM_24"
		} else {
			variant = "STANDARD"
		}
	}

	rawPrice, err := s.getCurrentPriceInTRY(req.Type)
	if err != nil || rawPrice <= 0 {
		return errors.New("piyasa verisi su an alinamiyor")
	}

	multiplier := 1.0
	if req.Type == "GOLD" {
		if val, ok := GoldFactors[variant]; ok {
			multiplier = val
		}
	}
	unitPrice := rawPrice * multiplier

	asset, err := s.repo.GetAsset(userID, req.Type, variant)
	if err != nil {
		asset = &model.Asset{UserID: userID, Type: req.Type, Variant: variant, Amount: 0}
	}

	if req.Action == "add" {
		asset.Amount += req.Amount
	} else if req.Action == "subtract" {
		if asset.Amount < req.Amount {
			return fmt.Errorf("yetersiz bakiye: %.2f adet/gr %s mevcut", asset.Amount, variant)
		}
		asset.Amount -= req.Amount
	}

	tDate := time.Now()
	if req.TransactionDate != nil {
		tDate = *req.TransactionDate
	}

	tx := &model.Transaction{
		Type:            req.Action,
		Amount:          req.Amount,
		Price:           unitPrice,
		TransactionDate: tDate,
	}

	return s.repo.UpdateWithLog(asset, tx)
}

func (s *AssetService) getCurrentPriceInTRY(assetType string) (float64, error) {
	rates, err := s.marketService.GetCurrencyRates()
	if err != nil && assetType != "BTC" && assetType != "GOLD" && assetType != "SILVER" {
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
	case "BTC":
		return s.marketService.GetCryptoPrice("BTC")
	case "TRY":
		return 1.0, nil
	default:
		return 1.0, nil
	}
}

// GÜNCELLENDİ: Artık baseVariant parametresi de alıyor
func (s *AssetService) GetPortfolioSummary(userID int64, baseCurrency string, baseVariant string) (*dto.PortfolioResponse, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 1. Hedef para biriminin TRY değeri
	basePriceTRY, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePriceTRY <= 0 {
		basePriceTRY = 1
	}

	// 2. Eğer hedef GOLD ise ve bir varyant seçildiyse (örn: Çeyrek), fiyatı ona göre ayarla
	targetMultiplier := 1.0
	if baseCurrency == "GOLD" && baseVariant != "" && baseVariant != "STANDARD" {
		if val, ok := GoldFactors[baseVariant]; ok {
			targetMultiplier = val
		}
	}
	// Hedef birim fiyatı (Örn: 1 Çeyrek Altın kaç TL?)
	finalBasePrice := basePriceTRY * targetMultiplier

	var details []dto.AssetResponse
	var totalValue float64

	for _, a := range assets {
		// Varlığın TRY bazındaki ham fiyatı
		rawPriceTRY, err := s.getCurrentPriceInTRY(a.Type)
		if err != nil {
			continue
		}

		// Varlığın kendi çarpanı (Örn: Elimdeki 22 Ayar Gram'ın çarpanı)
		multiplier := 1.0
		if a.Type == "GOLD" {
			if val, ok := GoldFactors[a.Variant]; ok {
				multiplier = val
			}
		}

		// Birim fiyatı hedef birime çevir
		currentUnitPrice := (rawPriceTRY * multiplier) / finalBasePrice

		totalAssetValue := a.Amount * currentUnitPrice
		totalValue += totalAssetValue

		details = append(details, dto.AssetResponse{
			ID:           a.ID,
			Type:         a.Type,
			Variant:      a.Variant,
			Amount:       a.Amount,
			CurrentPrice: currentUnitPrice,
			ValueInBase:  totalAssetValue,
		})
	}

	if totalValue > 0 {
		for i := range details {
			details[i].Allocation = (details[i].ValueInBase / totalValue) * 100
		}
	}

	return &dto.PortfolioResponse{
		Assets:     details,
		TotalValue: totalValue,
		BaseAsset:  baseCurrency,
	}, nil
}

func (s *AssetService) GetUserTransactionsWithCurrency(userID int64, baseCurrency string) ([]dto.TransactionResponse, error) {
	txs, err := s.repo.GetTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	basePrice, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePrice <= 0 {
		basePrice = 1
	}

	var response []dto.TransactionResponse
	for _, tx := range txs {
		if tx.Asset == nil {
			continue
		}
		response = append(response, dto.TransactionResponse{
			ID:              tx.ID,
			AssetID:         tx.AssetID,
			Type:            tx.Type,
			AssetType:       tx.Asset.Type,
			Variant:         tx.Asset.Variant,
			Amount:          tx.Amount,
			Price:           tx.Price / basePrice,
			TransactionDate: tx.TransactionDate.Add(3 * time.Hour),
			CreatedAt:       tx.CreatedAt,
		})
	}
	return response, nil
}

func (s *AssetService) GenerateTransactionReceipt(tx *model.Transaction, baseCurrency string) ([]byte, error) {
	if tx.Asset == nil {
		return nil, errors.New("islem detayi eksik")
	}
	basePrice, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePrice <= 0 {
		basePrice = 1
	}
	convertedPrice := tx.Price / basePrice
	userName, _ := s.repo.GetUserName(tx.Asset.UserID)

	txTypeTr := "Ekleme"
	if tx.Type == "subtract" {
		txTypeTr = "Cikarma"
	}

	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, s.tr("FINTRACK PRO - ISLEM DEKONTU"))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Varlik: %s", tx.Asset.Type)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Tur/Ayar: %s", tx.Asset.Variant)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Islem Tipi: %s", txTypeTr)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Miktar: %.2f", tx.Amount)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Birim Fiyat: %.2f %s", convertedPrice, baseCurrency)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Toplam Tutar: %.2f %s", tx.Amount*convertedPrice, baseCurrency)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Islem Tarihi: %s", tx.TransactionDate.Add(3*time.Hour).Format("02.01.2006 15:04"))))
	pdf.Ln(20)

	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Bu belge %s tarafindan uretilmistir.", userName)))

	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateFullPortfolioReceipt(userID int64, baseCurrency string) ([]byte, error) {
	// Not: Rapor şimdilik standart baseCurrency üzerinden devam ediyor
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, "STANDARD")
	if err != nil {
		return nil, err
	}
	userName, _ := s.repo.GetUserName(userID)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 15, s.tr("FINTRACK PRO - PORTFOY OZETI"))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, s.tr("Varlik"))
	pdf.Cell(40, 10, s.tr("Varyant"))
	pdf.Cell(30, 10, s.tr("Miktar"))
	pdf.Cell(40, 10, s.tr("Birim Deger"))
	pdf.Cell(40, 10, s.tr(fmt.Sprintf("Toplam (%s)", baseCurrency)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for _, a := range summary.Assets {
		pdf.Cell(40, 8, a.Type)
		pdf.Cell(40, 8, a.Variant)
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", a.Amount))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.CurrentPrice))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.ValueInBase))
		pdf.Ln(8)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("TOPLAM PORTFOY DEGERI: %.2f %s", summary.TotalValue, baseCurrency)))
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Kullanici: %s", userName)))

	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateExcelReport(userID int64, baseCurrency string) ([]byte, error) {
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, "STANDARD")
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetName := "Portfoy"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	f.SetCellValue(sheetName, "A1", "Varlik")
	f.SetCellValue(sheetName, "B1", "Varyant")
	f.SetCellValue(sheetName, "C1", "Miktar")
	f.SetCellValue(sheetName, "D1", fmt.Sprintf("Birim Fiyat (%s)", baseCurrency))
	f.SetCellValue(sheetName, "E1", fmt.Sprintf("Toplam (%s)", baseCurrency))

	row := 2
	for _, a := range summary.Assets {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), a.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), a.Variant)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), a.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), a.CurrentPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), a.ValueInBase)
		row++
	}

	f.SetCellValue(sheetName, fmt.Sprintf("D%d", row+1), "GENEL TOPLAM")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", row+1), summary.TotalValue)

	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}

func (s *AssetService) GetUserTransactions(userID int64) ([]model.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}
func (s *AssetService) GetTransactionByID(userID int64, txID int64) (*model.Transaction, error) {
	return s.repo.GetTransactionByID(txID, userID)
}
