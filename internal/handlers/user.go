package handlers

import (
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
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}
	if err := h.userService.Register(req.Email, req.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kayıt olunamadı"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Kayıt başarılı!"})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}
	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Giriş başarısız"})
	}
	return c.JSON(fiber.Map{"token": token})
}
