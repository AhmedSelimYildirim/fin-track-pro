package repository

import (
	"context"
	"fin-track-pro/internal/core/models"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(context.Background())
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}
