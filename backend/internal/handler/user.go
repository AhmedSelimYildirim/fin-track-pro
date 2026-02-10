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
		return c.Status(400).JSON(fiber.Map{"error": "Gecersiz veri formati."})
	}
	if err := h.userService.Register(req.Username, req.Email, req.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Kayit basarili! FinTrack Pro'ya hos geldiniz."})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gecersiz veri formati."})
	}
	token, user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"token":    token,
		"username": user.Username,
		"email":    user.Email,
	})
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
		return c.Status(401).JSON(fiber.Map{"error": "Yetkisiz erisim."})
	}
	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gecersiz veri."})
	}
	if err := h.userService.Update(userID, req.Username, req.Email, req.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Profil guncellendi!"})
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
		return c.Status(401).JSON(fiber.Map{"error": "Yetkisiz erisim."})
	}
	if err := h.userService.Delete(userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Silme islemi basarisiz."})
	}
	return c.JSON(fiber.Map{"message": "Hesabin basariyla silindi. Hosca kal!"})
}
