package service

import (
	"context"
	"fin-track-pro/internal/dto"
	"fin-track-pro/internal/model"
	"fin-track-pro/internal/repository"
)

type CalendarService struct {
	calendarRepo *repository.CalendarRepository
}

func NewCalendarService(repo *repository.CalendarRepository) *CalendarService {
	return &CalendarService{calendarRepo: repo}
}

func (s *CalendarService) CreateReminder(userID int64, req dto.CreateReminderRequest) error {
	reminder := &model.Reminder{
		UserID:     userID,
		Title:      req.Title,
		TargetDate: req.TargetDate,
	}
	return s.calendarRepo.SaveReminder(context.Background(), reminder)
}

func (s *CalendarService) GetUserReminders(userID int64) ([]model.Reminder, error) {
	return s.calendarRepo.GetRemindersByUserID(context.Background(), userID)
}
