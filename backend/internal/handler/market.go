package handler

import (
	"fin-track-pro/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MarketHandler struct {
	marketService *service.MarketService
}

func NewMarketHandler(s *service.MarketService) *MarketHandler {
	return &MarketHandler{marketService: s}
}

func (h *MarketHandler) GetRates(c *fiber.Ctx) error {
	price, err := h.marketService.GetMetalPrice("XAU")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"gold_gram_try": price})
}
