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

type AssetService struct {
	repo          *repository.AssetRepository
	marketService *MarketService
}

func NewAssetService(repo *repository.AssetRepository, market *MarketService) *AssetService {
	return &AssetService{repo: repo, marketService: market}
}

func (s *AssetService) tr(text string) string {
	replacer := strings.NewReplacer(
		"ğ", "g", "Ğ", "G",
		"ü", "u", "Ü", "U",
		"ş", "s", "Ş", "S",
		"ı", "i", "İ", "I",
		"ö", "o", "Ö", "O",
		"ç", "c", "Ç", "C",
	)
	return replacer.Replace(text)
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
		return errors.New("piyasa verisi su an alinamiyor")
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
			return fmt.Errorf("yetersiz bakiye: %s kasasinda sadece %.4f var", req.Type, asset.Amount)
		}
		if asset.Amount > 0 {
			asset.TotalCost -= req.Amount * (asset.TotalCost / asset.Amount)
		}
		asset.Amount -= req.Amount
	}
	tDate := time.Now()
	if req.TransactionDate != nil {
		tDate = *req.TransactionDate
	}
	tx := &model.Transaction{
		UserID:          int64(userID),
		Type:            req.Action,
		AssetType:       req.Type,
		Amount:          req.Amount,
		Price:           unitPrice,
		Ayar:            ayar,
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
	for _, a := range assets {
		rawPriceTRY, err := s.getCurrentPriceInTRY(a.Type)
		if err != nil || rawPriceTRY <= 0 {
			continue
		}
		displayAyar := 0
		if a.Type == "GOLD" {
			displayAyar = a.Ayar
			rawPriceTRY *= (float64(a.Ayar) / 24.0)
		}
		rate := rawPriceTRY / basePriceInTRY
		val := a.Amount * rate
		totalValue += val
		details = append(details, dto.AssetResponse{
			Type:         a.Type,
			Amount:       a.Amount,
			Ayar:         displayAyar,
			CurrentPrice: rate,
			ValueInBase:  val,
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

func (s *AssetService) GetUserTransactionsWithCurrency(userID uint, baseCurrency string, targetAyar int) ([]dto.TransactionResponse, error) {
	txs, err := s.repo.GetTransactionsByUserID(userID)
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
	var response []dto.TransactionResponse
	for _, tx := range txs {
		displayAyar := 0
		if tx.AssetType == "GOLD" {
			displayAyar = tx.Ayar
		}
		response = append(response, dto.TransactionResponse{
			ID:              tx.ID,
			Type:            tx.Type,
			AssetType:       tx.AssetType,
			Amount:          tx.Amount,
			Price:           tx.Price / basePriceInTRY,
			Ayar:            displayAyar,
			TransactionDate: tx.TransactionDate,
			CreatedAt:       tx.CreatedAt,
		})
	}
	return response, nil
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
	pdf.Cell(0, 10, s.tr("FINTRACK PRO - ISLEM DEKONTU"))
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Varlik: %s", tx.AssetType)))
	pdf.Ln(8)
	if tx.AssetType == "GOLD" && tx.Ayar > 0 {
		pdf.Cell(0, 10, s.tr(fmt.Sprintf("Ayar: %d Ayar", tx.Ayar)))
		pdf.Ln(8)
	}
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Islem Tipi: %s", txTypeTr)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Miktar: %.4f", tx.Amount)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Birim Fiyat: "+priceFormat+" %s", convertedPrice, baseCurrency)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Toplam Tutar: %.4f %s", tx.Amount*convertedPrice, baseCurrency)))
	pdf.Ln(8)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Islem Tarihi: %s", tx.TransactionDate.Format("02.01.2006"))))
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Bu belge %s tarafindan uretilmistir.", userName)))
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
	pdf.Cell(0, 15, s.tr("FINTRACK PRO - GENEL PORTFOY DEKONTU"))
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, s.tr("Varlik"))
	pdf.Cell(40, 10, s.tr("Miktar"))
	pdf.Cell(50, 10, s.tr("Birim Deger"))
	pdf.Cell(50, 10, s.tr(fmt.Sprintf("Toplam (%s)", baseCurrency)))
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
		pdf.Cell(40, 8, s.tr(vName))
		pdf.Cell(40, 8, fmt.Sprintf("%.4f", a.Amount))
		pdf.Cell(50, 8, fmt.Sprintf(pFormat, a.CurrentPrice))
		pdf.Cell(50, 8, fmt.Sprintf("%.4f", a.ValueInBase))
		pdf.Ln(8)
	}
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("TOPLAM PORTFOY DEGERI: %.4f %s", summary.TotalValue, baseCurrency)))
	pdf.Ln(15)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, s.tr(fmt.Sprintf("Rapor Tarihi: %s | Kullanici: %s", time.Now().Format("02.01.2006 15:04"), userName)))
	var buf bytes.Buffer
	pdf.Output(&buf)
	return buf.Bytes(), nil
}

func (s *AssetService) GenerateExcelReport(userID uint) ([]byte, error) {
	txs, err := s.repo.GetTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}
	f := excelize.NewFile()
	sheetName := "Portföy Raporu"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	styleTitle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1F4E78"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#2E75B6"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border:    []excelize.Border{{Type: "left", Color: "FFFFFF", Style: 1}, {Type: "right", Color: "FFFFFF", Style: 1}},
	})

	styleData, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border:    []excelize.Border{{Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}},
	})

	f.MergeCell(sheetName, "A1", "D1")
	f.SetCellValue(sheetName, "A1", "VARLIK ÖZETİ")
	f.SetCellStyle(sheetName, "A1", "D1", styleTitle)

	summaryHeaders := []string{"Varlık Tipi", "Toplam Miktar", "Ortalama Maliyet (Tahmini)", "Son İşlem Tarihi"}
	for i, h := range summaryHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, h)
	}
	f.SetCellStyle(sheetName, "A2", "D2", styleHeader)

	assetsMap := make(map[string]float64)
	lastDateMap := make(map[string]time.Time)

	for _, tx := range txs {
		key := tx.AssetType
		if tx.AssetType == "GOLD" && tx.Ayar > 0 {
			key = fmt.Sprintf("GOLD (%dK)", tx.Ayar)
		}

		if tx.Type == "add" {
			assetsMap[key] += tx.Amount
		} else {
			assetsMap[key] -= tx.Amount
		}

		if tx.TransactionDate.After(lastDateMap[key]) {
			lastDateMap[key] = tx.TransactionDate
		}
	}

	row := 3
	for k, v := range assetsMap {
		if v <= 0 {
			continue
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), k)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), v)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), lastDateMap[k].Format("02.01.2006"))
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("D%d", row), styleData)
		row++
	}

	startRow := row + 3
	f.MergeCell(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("H%d", startRow))
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "DETAYLI İŞLEM GEÇMİŞİ")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("H%d", startRow), styleTitle)

	txHeaders := []string{"ID", "Tarih", "İşlem Tipi", "Varlık", "Ayar", "Miktar", "Birim Fiyat (TRY)", "Toplam Tutar (TRY)"}
	headerRow := startRow + 1
	for i, h := range txHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheetName, cell, h)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", headerRow), fmt.Sprintf("H%d", headerRow), styleHeader)

	dataRow := headerRow + 1
	for _, tx := range txs {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", dataRow), tx.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", dataRow), tx.TransactionDate.Format("02.01.2006 15:04"))

		typeTr := "Ekleme"
		if tx.Type == "subtract" {
			typeTr = "Çıkarma"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", dataRow), typeTr)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", dataRow), tx.AssetType)

		if tx.Ayar > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", dataRow), fmt.Sprintf("%d Ayar", tx.Ayar))
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", dataRow), "-")
		}

		f.SetCellValue(sheetName, fmt.Sprintf("F%d", dataRow), tx.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", dataRow), tx.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", dataRow), tx.Amount*tx.Price)

		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", dataRow), fmt.Sprintf("H%d", dataRow), styleData)
		dataRow++
	}

	f.SetColWidth(sheetName, "A", "A", 10)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 15)
	f.SetColWidth(sheetName, "E", "E", 15)
	f.SetColWidth(sheetName, "F", "F", 15)
	f.SetColWidth(sheetName, "G", "H", 20)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *AssetService) GetUserTransactions(userID uint) ([]model.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

func (s *AssetService) GetTransactionByID(userID uint, txID int64) (*model.Transaction, error) {
	return s.repo.GetTransactionByID(txID, int64(userID))
}
