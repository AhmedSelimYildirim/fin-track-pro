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

func (s *AssetService) GetPortfolioSummary(userID int64, baseCurrency string, baseVariant string) (*dto.PortfolioResponse, error) {
	assets, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	basePriceTRY, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePriceTRY <= 0 {
		basePriceTRY = 1
	}

	targetMultiplier := 1.0
	if baseCurrency == "GOLD" && baseVariant != "" && baseVariant != "STANDARD" {
		if val, ok := GoldFactors[baseVariant]; ok {
			targetMultiplier = val
		}
	}
	finalBasePrice := basePriceTRY * targetMultiplier

	var details []dto.AssetResponse
	var totalValue float64

	for _, a := range assets {
		rawPriceTRY, err := s.getCurrentPriceInTRY(a.Type)
		if err != nil {
			continue
		}

		multiplier := 1.0
		if a.Type == "GOLD" {
			if val, ok := GoldFactors[a.Variant]; ok {
				multiplier = val
			}
		}

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
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Islem Tarihi: %s", tx.TransactionDate.Add(3*time.Hour).Format("02.01.2006"))))
	pdf.Ln(20)

	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Bu belge %s tarafindan uretilmistir.", userName)))

	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateFullPortfolioReceipt(userID int64, baseCurrency string, baseVariant string, filterType string) ([]byte, error) {
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, baseVariant)
	if err != nil {
		return nil, err
	}
	userName, _ := s.repo.GetUserName(userID)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	title := "FINTRACK PRO - PORTFOY OZETI"
	if filterType != "" {
		title = fmt.Sprintf("FINTRACK PRO - %s PORTFOYU", filterType)
	}

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 15, s.tr(title))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, s.tr("Varlik"))
	pdf.Cell(40, 10, s.tr("Varyant"))
	pdf.Cell(30, 10, s.tr("Miktar"))
	pdf.Cell(40, 10, s.tr("Birim Deger"))
	unitLabel := baseCurrency
	if baseVariant != "" && baseVariant != "STANDARD" {
		unitLabel = fmt.Sprintf("%s (%s)", baseCurrency, baseVariant)
	}
	pdf.Cell(40, 10, s.tr(fmt.Sprintf("Toplam (%s)", unitLabel)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	filteredTotal := 0.0

	for _, a := range summary.Assets {
		if filterType != "" && a.Type != filterType {
			continue
		}
		pdf.Cell(40, 8, a.Type)
		pdf.Cell(40, 8, a.Variant)
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", a.Amount))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.CurrentPrice))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.ValueInBase))
		pdf.Ln(8)
		filteredTotal += a.ValueInBase
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("TOPLAM DEGER: %.2f %s", filteredTotal, unitLabel)))
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Kullanici: %s", userName)))

	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateExcelReport(userID int64, baseCurrency string, baseVariant string, filterType string) ([]byte, error) {
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, baseVariant)
	if err != nil {
		return nil, err
	}
	txs, _ := s.repo.GetTransactionsByUserID(userID)

	f := excelize.NewFile()
	sheetName := "Varliklar"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	unitLabel := baseCurrency
	if baseVariant != "" && baseVariant != "STANDARD" {
		unitLabel = fmt.Sprintf("%s (%s)", baseCurrency, baseVariant)
	}

	f.SetCellValue(sheetName, "A1", "Varlik")
	f.SetCellValue(sheetName, "B1", "Varyant")
	f.SetCellValue(sheetName, "C1", "Miktar")
	f.SetCellValue(sheetName, "D1", fmt.Sprintf("Birim Fiyat (%s)", unitLabel))
	f.SetCellValue(sheetName, "E1", fmt.Sprintf("Toplam (%s)", unitLabel))

	row := 2
	filteredTotal := 0.0
	for _, a := range summary.Assets {
		if filterType != "" && a.Type != filterType {
			continue
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), a.Type)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), a.Variant)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), a.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), a.CurrentPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), a.ValueInBase)
		row++
		filteredTotal += a.ValueInBase
	}
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", row+1), "TOPLAM")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", row+1), filteredTotal)

	// İŞLEM GEÇMİŞİ SAYFASI
	histSheet := "Islem Gecmisi"
	f.NewSheet(histSheet)
	f.SetCellValue(histSheet, "A1", "Tarih")
	f.SetCellValue(histSheet, "B1", "Islem")
	f.SetCellValue(histSheet, "C1", "Varlik")
	f.SetCellValue(histSheet, "D1", "Miktar")
	f.SetCellValue(histSheet, "E1", "Birim Fiyat (Islem Aninda)")

	hRow := 2
	for _, tx := range txs {
		if tx.Asset == nil {
			continue
		}
		if filterType != "" && tx.Asset.Type != filterType {
			continue
		}

		f.SetCellValue(histSheet, fmt.Sprintf("A%d", hRow), tx.TransactionDate.Format("02.01.2006"))
		f.SetCellValue(histSheet, fmt.Sprintf("B%d", hRow), tx.Type)
		f.SetCellValue(histSheet, fmt.Sprintf("C%d", hRow), fmt.Sprintf("%s - %s", tx.Asset.Type, tx.Asset.Variant))
		f.SetCellValue(histSheet, fmt.Sprintf("D%d", hRow), tx.Amount)
		f.SetCellValue(histSheet, fmt.Sprintf("E%d", hRow), tx.Price)
		hRow++
	}

	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}

func (s *AssetService) GetUserTransactions(userID int64) ([]model.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}
func (s *AssetService) GetTransactionByID(userID int64, txID int64) (*model.Transaction, error) {
	return s.repo.GetTransactionByID(txID, userID)
}
