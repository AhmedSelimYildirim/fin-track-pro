package dto

import "time"

type CreateReminderRequest struct {
	Title      string    `json:"title" validate:"required"`
	TargetDate time.Time `json:"target_date" validate:"required"`
}

type ReminderResponse struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	TargetDate time.Time `json:"target_date"`
	CreatedAt  time.Time `json:"created_at"`
}
