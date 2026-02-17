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
	tx := &model.Transaction{Type: req.Action, Amount: req.Amount, Price: unitPrice, TransactionDate: tDate}
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
		rawPriceTRY, _ := s.getCurrentPriceInTRY(a.Type)
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
			ID: a.ID, Type: a.Type, Variant: a.Variant, Amount: a.Amount,
			CurrentPrice: currentUnitPrice, ValueInBase: totalAssetValue,
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
	pdf.Cell(0, 10, s.tr("İŞLEM DEKONTU"))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Varlık: %s (%s)", tx.Asset.Type, tx.Asset.Variant)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Miktar: %.2f", tx.Amount)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Birim Fiyat: %.2f %s", convertedPrice, baseCurrency)))
	pdf.Ln(8)
	trDate := tx.TransactionDate.Add(3 * time.Hour)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Tarih: %s", trDate.Format("02.01.2006"))))
	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Müşteri: %s", userName)))
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

	title := "GENEL PORTFOY RAPORU"
	if filterType != "" {
		title = fmt.Sprintf("%s PORTFOY RAPORU", filterType)
	}
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 15, s.tr(title))
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, s.tr("Varlık"))
	pdf.Cell(40, 10, s.tr("Varyant"))
	pdf.Cell(30, 10, s.tr("Miktar"))
	pdf.Cell(40, 10, s.tr(fmt.Sprintf("Değer (%s)", baseCurrency)))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	var filteredTotal float64
	for _, a := range summary.Assets {
		if filterType != "" && a.Type != filterType {
			continue
		}
		pdf.Cell(40, 8, a.Type)
		pdf.Cell(40, 8, a.Variant)
		pdf.Cell(30, 8, fmt.Sprintf("%.2f", a.Amount))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f", a.ValueInBase))
		pdf.Ln(8)
		filteredTotal += a.ValueInBase
	}
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("TOPLAM DEĞER: %.2f %s", filteredTotal, baseCurrency)))
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Müşteri: %s", userName)))
	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateExcelReport(userID int64, baseCurrency string, baseVariant string, filterType string) ([]byte, error) {
	summary, _ := s.GetPortfolioSummary(userID, baseCurrency, baseVariant)
	txs, _ := s.repo.GetTransactionsByUserID(userID)
	f := excelize.NewFile()

	sheet := "Portfoy Özeti"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1E293B"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	reportTitle := "GENEL VARLIK LİSTESİ"
	if filterType != "" {
		reportTitle = fmt.Sprintf("%s VARLIK LİSTESİ", filterType)
	}
	f.SetCellValue(sheet, "A1", s.tr(reportTitle))
	f.SetCellValue(sheet, "A3", "Varlık")
	f.SetCellValue(sheet, "B3", "Varyant")
	f.SetCellValue(sheet, "C3", "Miktar")
	f.SetCellValue(sheet, "D3", s.tr(fmt.Sprintf("Birim Fiyat (%s)", baseCurrency)))
	f.SetCellValue(sheet, "E3", s.tr(fmt.Sprintf("Toplam (%s)", baseCurrency)))
	f.SetCellStyle(sheet, "A3", "E3", headerStyle)

	row := 4
	var total float64
	for _, a := range summary.Assets {
		if filterType != "" && a.Type != filterType {
			continue
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), a.Type)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), a.Variant)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), a.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), a.CurrentPrice)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), a.ValueInBase)
		total += a.ValueInBase
		row++
	}
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "GENEL TOPLAM:")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), total)
	f.SetColWidth(sheet, "A", "E", 18)

	sheet2 := "Islem Gecmisi"
	f.NewSheet(sheet2)

	f.SetCellValue(sheet2, "A1", s.tr("TUM İSLEMLER DÖKÜMÜ"))
	f.SetCellValue(sheet2, "A3", "Tarih")
	f.SetCellValue(sheet2, "B3", "İşlem")
	f.SetCellValue(sheet2, "C3", "Varlık")
	f.SetCellValue(sheet2, "D3", "Varyant")
	f.SetCellValue(sheet2, "E3", "Miktar")
	f.SetCellValue(sheet2, "F3", s.tr(fmt.Sprintf("İşlem Değeri (%s)", baseCurrency)))
	f.SetCellStyle(sheet2, "A3", "F3", headerStyle)

	basePrice, _ := s.getCurrentPriceInTRY(baseCurrency)
	if basePrice <= 0 {
		basePrice = 1
	}

	row = 4
	for i := len(txs) - 1; i >= 0; i-- {
		t := txs[i]

		if filterType != "" && t.Asset.Type != filterType {
			continue
		}

		trDate := t.TransactionDate.Add(3 * time.Hour)
		op := "Ekleme"
		if t.Type == "subtract" {
			op = "Cikarma"
		}

		// O anki işlem fiyatını seçilen para birimine çevir
		convertedTxPrice := t.Price / basePrice

		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), trDate.Format("02.01.2006"))
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), op)
		f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), t.Asset.Type)
		f.SetCellValue(sheet2, fmt.Sprintf("D%d", row), t.Asset.Variant)
		f.SetCellValue(sheet2, fmt.Sprintf("E%d", row), t.Amount)
		f.SetCellValue(sheet2, fmt.Sprintf("F%d", row), t.Amount*convertedTxPrice)
		row++
	}
	f.SetColWidth(sheet2, "A", "F", 18)

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
