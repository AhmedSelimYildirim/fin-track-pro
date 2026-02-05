package handlers

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
	rates, err := h.marketService.GetCurrencyRates()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	btc, _ := h.marketService.GetCryptoPrice("bitcoin")
	gold, _ := h.marketService.GetMetalPrice("XAU")

	return c.JSON(fiber.Map{
		"currency": rates,
		"bitcoin":  btc,
		"gold":     gold,
	})
}
