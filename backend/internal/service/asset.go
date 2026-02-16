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
		variant = "GRAM_24"
		if req.Type != "GOLD" {
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
			return fmt.Errorf("yetersiz bakiye")
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
	default:
		return 1.0, nil
	}
}

func (s *AssetService) GetPortfolioSummary(userID int64, baseCurrency string, baseVariant string) (*dto.PortfolioResponse, error) {
	assets, _ := s.repo.GetByUserID(userID)
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
	return &dto.PortfolioResponse{Assets: details, TotalValue: totalValue, BaseAsset: baseCurrency}, nil
}

func (s *AssetService) GenerateTransactionReceipt(tx *model.Transaction, baseCurrency string) ([]byte, error) {
	basePrice, _ := s.getCurrentPriceInTRY(baseCurrency)
	convertedPrice := tx.Price / basePrice
	userName, _ := s.repo.GetUserName(tx.UserID)
	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, s.tr("ISLEM DEKONTU"))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Varlik: %s (%s)", tx.Asset.Type, tx.Asset.Variant)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Miktar: %.2f", tx.Amount)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Birim Fiyat: %.2f %s", convertedPrice, baseCurrency)))
	pdf.Ln(8)
	trDate := tx.TransactionDate.Add(3 * time.Hour)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Tarih: %s", trDate.Format("02.01.2006"))))
	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Musteri: %s", userName)))
	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateFullPortfolioReceipt(userID int64, baseCurrency string, baseVariant string, filterType string) ([]byte, error) {
	summary, err := s.GetPortfolioSummary(userID, baseCurrency, baseVariant)
	if err != nil {
		return nil, err
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 15, s.tr(fmt.Sprintf("%s PORTFOY RAPORU", filterType)))
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, s.tr("Varlik"))
	pdf.Cell(30, 10, s.tr("Miktar"))
	pdf.Cell(40, 10, s.tr(fmt.Sprintf("Deger (%s)", baseCurrency)))
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	var filteredTotal float64
	for _, a := range summary.Assets {
		if filterType != "" && a.Type != filterType {
			continue
		}
		pdf.Cell(40, 8, s.tr(fmt.Sprintf("%s-%s", a.Type, a.Variant)))
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", a.Amount))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.ValueInBase))
		pdf.Ln(8)
		filteredTotal += a.ValueInBase
	}
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("GENEL TOPLAM: %.2f %s", filteredTotal, baseCurrency)))
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
	sheet := "Varlik Ekstresi"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1E293B"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	f.SetCellValue(sheet, "A1", s.tr(fmt.Sprintf("PORTFOY: %s", filterType)))
	f.SetCellValue(sheet, "A2", "Varlik Listesi")
	f.SetCellValue(sheet, "A3", "Varlik")
	f.SetCellValue(sheet, "B3", "Miktar")
	f.SetCellValue(sheet, "C3", s.tr(fmt.Sprintf("Toplam Deger (%s)", baseCurrency)))
	f.SetCellStyle(sheet, "A3", "C3", headerStyle)

	row := 4
	var total float64
	for _, a := range summary.Assets {
		if filterType != "" && a.Type != filterType {
			continue
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%s (%s)", a.Type, a.Variant))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), a.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), a.ValueInBase)
		total += a.ValueInBase
		row++
	}
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOPLAM:")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), total)

	row += 3
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Islem Gecmisi")
	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Tarih")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Islem")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Miktar")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.tr(fmt.Sprintf("Birim Fiyat (%s)", baseCurrency)))
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), headerStyle)
	row++

	basePrice, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePrice <= 0 {
		basePrice = 1
	}

	for i := len(txs) - 1; i >= 0; i-- {
		t := txs[i]
		if filterType != "" && t.Asset.Type != filterType {
			continue
		}
		trDate := t.TransactionDate.Add(3 * time.Hour)
		opType := "Ekleme"
		if t.Type == "subtract" {
			opType = "Cikarma"
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), trDate.Format("02.01.2006"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), opType)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Price/basePrice)
		row++
	}

	f.SetColWidth(sheet, "A", "D", 20)
	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}

func (s *AssetService) GetUserTransactionsWithCurrency(userID int64, baseCurrency string) ([]dto.TransactionResponse, error) {
	txs, _ := s.repo.GetTransactionsByUserID(userID)
	basePrice, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePrice <= 0 {
		basePrice = 1
	}
	var res []dto.TransactionResponse
	for _, tx := range txs {
		res = append(res, dto.TransactionResponse{
			ID: tx.ID, Type: tx.Type, AssetType: tx.Asset.Type, Variant: tx.Asset.Variant,
			Amount: tx.Amount, Price: tx.Price / basePrice,
			TransactionDate: tx.TransactionDate.Add(3 * time.Hour),
		})
	}
	return res, nil
}

func (s *AssetService) GetTransactionByID(userID, txID int64) (*model.Transaction, error) {
	return s.repo.GetTransactionByID(txID, userID)
}
