package handlers

import (
	"fin-track-pro/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CalendarHandler struct {
	calendarService *service.CalendarService
}

func NewCalendarHandler(s *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{calendarService: s}
}

func (h *CalendarHandler) AddEvent(c *fiber.Ctx) error {
	userID := int64(c.Locals("user_id").(uint))
	var req struct {
		Title      string    `json:"title"`
		TargetDate time.Time `json:"target_date"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gecersiz format"})
	}
	err := h.calendarService.CreateReminder(userID, req.Title, req.TargetDate)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Hatırlatıcı kaydedilemedi"})
	}
	return c.JSON(fiber.Map{"message": "Hatırlatıcı başarıyla oluşturuldu !"})
}

func (h *CalendarHandler) ListReminders(c *fiber.Ctx) error {
	userID := int64(c.Locals("user_id").(uint))
	reminders, err := h.calendarService.GetUserReminders(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reminders)
}
