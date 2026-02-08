package handler

import (
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/service"

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
	var req dto.CreateReminderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gecersiz format"})
	}
	err := h.calendarService.CreateReminder(userID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Not kaydedilemedi"})
	}
	return c.JSON(fiber.Map{"message": "Not basariyla kaydedildi !"})
}

func (h *CalendarHandler) ListReminders(c *fiber.Ctx) error {
	userID := int64(c.Locals("user_id").(uint))
	notes, err := h.calendarService.GetUserReminders(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(notes)
}
