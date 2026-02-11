package handler

import (
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/service"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	service *service.AssetService
}

func NewAssetHandler(s *service.AssetService) *AssetHandler {
	return &AssetHandler{service: s}
}

func (h *AssetHandler) UpdateBalance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req dto.AssetCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gecersiz format"})
	}
	if err := h.service.ManageBalance(userID, req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Bakiye guncellendi !"})
}

func (h *AssetHandler) GetSummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	currency := c.Query("currency")
	if currency == "" {
		currency = c.Get("X-Currency", "TRY")
	}
	ayarStr := c.Query("ayar")
	if ayarStr == "" {
		ayarStr = c.Get("X-Ayar", "0")
	}
	ayar, _ := strconv.Atoi(ayarStr)
	summary, err := h.service.GetPortfolioSummary(userID, currency, ayar)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(summary)
}

func (h *AssetHandler) GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	currency := c.Query("currency")
	if currency == "" {
		currency = c.Get("X-Currency", "TRY")
	}
	ayarStr := c.Query("ayar")
	if ayarStr == "" {
		ayarStr = c.Get("X-Ayar", "0")
	}
	ayar, _ := strconv.Atoi(ayarStr)
	txs, err := h.service.GetUserTransactionsWithCurrency(userID, currency, ayar)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(txs)
}

func (h *AssetHandler) GetReceipt(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	userID := c.Locals("user_id").(uint)
	currency := c.Query("currency")
	if currency == "" {
		currency = c.Get("X-Currency", "TRY")
	}
	ayarStr := c.Query("ayar")
	if ayarStr == "" {
		ayarStr = c.Get("X-Ayar", "0")
	}
	ayar, _ := strconv.Atoi(ayarStr)
	tx, err := h.service.GetTransactionByID(userID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Islem bulunamadi"})
	}
	pdfBytes, err := h.service.GenerateTransactionReceipt(tx, currency, ayar)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Dekont uretilemedi"})
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=dekont_%d.pdf", id))
	return c.Send(pdfBytes)
}

func (h *AssetHandler) GetFullPortfolioReceipt(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	currency := c.Query("currency")
	if currency == "" {
		currency = c.Get("X-Currency", "TRY")
	}
	ayarStr := c.Query("ayar")
	if ayarStr == "" {
		ayarStr = c.Get("X-Ayar", "0")
	}
	ayar, _ := strconv.Atoi(ayarStr)
	pdfBytes, err := h.service.GenerateFullPortfolioReceipt(userID, currency, ayar)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Portfoy raporu uretilemedi"})
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=toplam_portfoy.pdf")
	return c.Send(pdfBytes)
}

func (h *AssetHandler) GetExcel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	currency := c.Query("currency")
	if currency == "" {
		currency = c.Get("X-Currency", "TRY")
	}
	ayarStr := c.Query("ayar")
	if ayarStr == "" {
		ayarStr = c.Get("X-Ayar", "0")
	}
	ayar, _ := strconv.Atoi(ayarStr)
	excelBytes, err := h.service.GenerateExcelReport(userID, currency, ayar)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Excel raporu uretilemedi"})
	}
	fileName := fmt.Sprintf("FinTrack_Rapor_%s_%s.xlsx", currency, time.Now().Add(3*time.Hour).Format("20060102"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	return c.Send(excelBytes)
}
