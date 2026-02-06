package service

import (
	"context"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
	"log"
	"time"
)

type CalendarService struct {
	calendarRepo *repository.CalendarRepository
}

func NewCalendarService(repo *repository.CalendarRepository) *CalendarService {
	return &CalendarService{calendarRepo: repo}
}

func (s *CalendarService) CreateReminder(userID int64, title string, remindAt time.Time) error {
	reminder := &models.Reminder{
		UserID:     userID,
		Title:      title,
		TargetDate: remindAt,
		IsSent:     false,
	}
	return s.calendarRepo.SaveReminder(context.Background(), reminder)
}

func (s *CalendarService) GetUserReminders(userID int64) ([]models.Reminder, error) {
	return s.calendarRepo.GetRemindersByUserID(context.Background(), userID)
}

func (s *CalendarService) ProcessPendingReminders() {
	ctx := context.Background()
	reminders, err := s.calendarRepo.GetPendingReminders(ctx)
	if err != nil {
		log.Printf("Hatirlatici cekme hatasi: %v", err)
		return
	}

	for _, r := range reminders {
		log.Printf("Bildirim Gonderiliyor: %s (User ID: %d)", r.Title, r.UserID)
		_ = s.calendarRepo.MarkAsSent(ctx, r.ID)
	}
}
