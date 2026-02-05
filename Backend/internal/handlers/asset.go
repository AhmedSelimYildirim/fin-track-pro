package handlers

import (
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/service"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	service       *service.AssetService
	marketService *service.MarketService
}

func NewAssetHandler(s *service.AssetService, m *service.MarketService) *AssetHandler {
	return &AssetHandler{
		service:       s,
		marketService: m,
	}
}

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var req dto.AssetCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gecersiz veri"})
	}
	userID := c.Locals("user_id").(uint)
	asset := &models.Asset{
		UserID:       int64(userID),
		Type:         req.Type,
		Amount:       req.Amount,
		Cost:         req.Cost,
		Note:         req.Note,
		PurchaseDate: req.PurchaseDate,
	}
	if asset.PurchaseDate.IsZero() {
		asset.PurchaseDate = time.Now()
	}
	if err := h.service.AddAsset(asset); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "varlik eklenemedi"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Varlik kaydedildi!", "id": asset.ID})
}

func (h *AssetHandler) GetSummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	summary, err := h.service.GetPortfolioSummary(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Hesaplama motoru durdu!"})
	}
	return c.Status(200).JSON(summary)
}

func (h *AssetHandler) GetAssets(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	assets, err := h.service.GetUserAssets(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "varliklar getirilemedi"})
	}
	return c.JSON(assets)
}

func (h *AssetHandler) ExportExcel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	fileBytes, err := h.service.ExportToExcel(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Excel olusturulamadi"})
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=fintrack_portfoy.xlsx")
	return c.Send(fileBytes)
}

func (h *AssetHandler) GetReceipt(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	userID := c.Locals("user_id").(uint)
	asset, err := h.service.GetAssetByID(userID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Varlik bulunamadi"})
	}
	pdfBytes, err := h.service.GenerateReceipt(asset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Dekont uretilemedi"})
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=dekont_%d.pdf", id))
	return c.Send(pdfBytes)
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	userID := c.Locals("user_id").(uint)
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gecersiz miktar"})
	}
	if err := h.service.UpdateAssetAmount(userID, id, req.Amount); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Varlık miktarı güncellendi Ahmed Selim!"})
}

func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	userID := c.Locals("user_id").(uint)
	if err := h.service.RemoveAsset(userID, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Varlık başarıyla silindi."})
}
