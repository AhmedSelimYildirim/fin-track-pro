package handlers

import (
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	service *service.AssetService
}

func NewAssetHandler(s *service.AssetService) *AssetHandler {
	return &AssetHandler{service: s}
}

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	var asset models.Asset
	if err := c.BodyParser(&asset); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "geçersiz veri"})
	}

	// userID normalde JWT'den gelecek, şimdilik manuel test edebilirsin
	if err := h.service.AddAsset(&asset); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "varlık eklenemedi"})
	}

	return c.Status(201).JSON(asset)
}
