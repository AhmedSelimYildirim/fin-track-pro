package service

import (
	"errors"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user := &models.User{Email: email, Password: hashedPassword}
	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	// 1. Kullanıcıyı email ile bul
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("kullanıcı bulunamadı")
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("e-posta veya şifre hatalı")
	}

	return GenerateToken(user.ID)
}
