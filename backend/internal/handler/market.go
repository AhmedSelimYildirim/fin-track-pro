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
	rates, err := h.marketService.GetCurrencyRates()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Doviz kurlari alinamadi"})
	}

	goldPrice, _ := h.marketService.GetMetalPrice("GOLD")
	silverPrice, _ := h.marketService.GetMetalPrice("SILVER")
	btcPrice, _ := h.marketService.GetCryptoPrice("BTC")

	response := fiber.Map{
		"USD":    rates["USD"],
		"EUR":    rates["EUR"],
		"GOLD":   goldPrice,
		"SILVER": silverPrice,
		"BTC":    btcPrice,
		"TRY":    1.0,
	}

	return c.JSON(response)
}
