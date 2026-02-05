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
