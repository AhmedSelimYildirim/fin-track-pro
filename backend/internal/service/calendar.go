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

func (s *CalendarService) GetUserReminders(userID int64) ([]dto.ReminderResponse, error) {
	reminders, err := s.calendarRepo.GetRemindersByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var response []dto.ReminderResponse
	for _, r := range reminders {
		response = append(response, dto.ReminderResponse{
			ID:         r.ID,
			Title:      r.Title,
			TargetDate: r.TargetDate,
			CreatedAt:  r.CreatedAt,
		})
	}
	return response, nil
}

func (s *CalendarService) DeleteReminder(userID int64, reminderID int64) error {
	return s.calendarRepo.DeleteReminder(context.Background(), userID, reminderID)
}
