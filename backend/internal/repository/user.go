package repository

import (
	"context"
	"fin-track-pro/internal/model"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(context.Background())
	return err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(context.Background())
	return err
}

func (r *UserRepository) Delete(id uint) error {
	_, err := r.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id).Exec(context.Background())
	return err
}
