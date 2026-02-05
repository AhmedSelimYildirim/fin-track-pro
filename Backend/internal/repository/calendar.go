package repository

import (
	"context"
	"fin-track-pro/internal/core/models"

	"github.com/uptrace/bun"
)

type CalendarRepository struct {
	db *bun.DB
}

func NewCalendarRepository(db *bun.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

func (r *CalendarRepository) SaveReminder(ctx context.Context, reminder *models.Reminder) error {
	_, err := r.db.NewInsert().Model(reminder).Exec(ctx)
	return err
}

func (r *CalendarRepository) GetRemindersByUserID(ctx context.Context, userID int64) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.NewSelect().
		Model(&reminders).
		Where("user_id = ?", userID).
		Order("target_date ASC").
		Scan(ctx)
	return reminders, err
}
