package dto

import "time"

type CreateReminderRequest struct {
	Title string `json:"title" validate:"required"`
}

type ReminderResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
