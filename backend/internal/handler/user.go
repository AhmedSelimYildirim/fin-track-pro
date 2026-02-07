package handler

import (
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri formatı."})
	}

	if err := h.userService.Register(req.Username, req.Email, req.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Kayıt başarılı! FinTrack Pro'ya hoş geldiniz."})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri formatı."})
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	var userID uint
	switch v := val.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return c.Status(401).JSON(fiber.Map{"error": "Yetkisiz erişim."})
	}

	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz veri."})
	}
	if err := h.userService.Update(userID, req.Username, req.Email); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Güncelleme başarısız."})
	}
	return c.JSON(fiber.Map{"message": "Profil güncellendi!"})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	var userID uint
	switch v := val.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return c.Status(401).JSON(fiber.Map{"error": "Yetkisiz erişim."})
	}

	if err := h.userService.Delete(userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Silme işlemi başarısız."})
	}
	return c.JSON(fiber.Map{"message": "Hesabın başarıyla silindi. Hoşça kal!"})
}
