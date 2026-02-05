package handlers

import (
	"fin-track-pro/internal/service"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	service *service.AssetService
}

func NewAssetHandler(s *service.AssetService, m *service.MarketService) *AssetHandler {
	return &AssetHandler{service: s}
}

func (h *AssetHandler) UpdateBalance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
		Action string  `json:"action"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gecersiz format"})
	}
	if err := h.service.ManageBalance(userID, req.Type, req.Amount, req.Action); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Bakiye guncellendi !"})
}

func (h *AssetHandler) GetSummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	summary, err := h.service.GetPortfolioSummary(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(summary)
}

func (h *AssetHandler) GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	txs, err := h.service.GetUserTransactions(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(txs)
}

func (h *AssetHandler) GetReceipt(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	userID := c.Locals("user_id").(uint)
	tx, err := h.service.GetTransactionByID(userID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Islem bulunamadi"})
	}
	pdfBytes, err := h.service.GenerateTransactionReceipt(tx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Dekont uretilemedi"})
	}
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=dekont_%d.pdf", id))
	return c.Send(pdfBytes)
}
