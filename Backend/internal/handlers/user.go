package handlers

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
		return c.Status(400).JSON(fiber.Map{"error": "geçersiz format"})
	}

	if err := h.userService.Register(req.Username, req.Email, req.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "kayıt başarısız"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Kayıt başarılı Ahmed Selim! 🚀"})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "geçersiz format"})
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}
