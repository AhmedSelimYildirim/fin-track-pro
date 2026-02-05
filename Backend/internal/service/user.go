package service

import (
	"errors"
	"fin-track-pro/internal/core/models"
	"fin-track-pro/internal/repository"
	"fin-track-pro/internal/utils"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, secret string) *UserService {
	return &UserService{repo: repo, jwtSecret: secret}
}

func (s *UserService) Register(username, email, password string) error {
	// utils paketindeki fonksiyonu kullanıyoruz
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("kullanıcı bulunamadı")
	}
	// utils paketindeki fonksiyonu kullanıyoruz
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("hatalı şifre")
	}
	// utils paketindeki fonksiyonu kullanıyoruz
	return utils.GenerateToken(uint(user.ID), s.jwtSecret)
}
