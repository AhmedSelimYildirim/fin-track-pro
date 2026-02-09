package repository

import (
	"context"
	"fin-track-pro/internal/model"

	"github.com/uptrace/bun"
)

type CalendarRepository struct {
	db *bun.DB
}

func NewCalendarRepository(db *bun.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

func (r *CalendarRepository) SaveReminder(ctx context.Context, reminder *model.Reminder) error {
	_, err := r.db.NewInsert().Model(reminder).Exec(ctx)
	return err
}

func (r *CalendarRepository) GetRemindersByUserID(ctx context.Context, userID int64) ([]model.Reminder, error) {
	var reminders []model.Reminder
	err := r.db.NewSelect().
		Model(&reminders).
		Where("user_id = ?", userID).
		Order("id DESC").
		Scan(ctx)
	return reminders, err
}

func (r *CalendarRepository) DeleteReminder(ctx context.Context, userID int64, reminderID int64) error {
	_, err := r.db.NewDelete().
		Model((*model.Reminder)(nil)).
		Where("id = ? AND user_id = ?", reminderID, userID).
		Exec(ctx)
	return err
}
