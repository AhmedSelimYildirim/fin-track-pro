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

// Basit bir test endpoint'i
func (h *CalendarHandler) AddEvent(c *fiber.Ctx) error {
	// Normalde buraya body'den veri gelir ama şimdilik manuel tetikleyelim
	err := h.calendarService.CreateReminder(1, "Altın Fiyatına Bak", time.Now().Add(24*time.Hour))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Hatırlatıcı kaydedilemedi"})
	}
	return c.JSON(fiber.Map{"message": "Hatırlatıcı başarıyla oluşturuldu Ahmed Selim!"})
}
