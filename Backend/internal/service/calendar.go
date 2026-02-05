package service

import (
	"context"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
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
