package service

import (
	"errors"
	"fin-track-pro/internal/model"
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
	if !utils.IsValidEmail(email) {
		return errors.New("Geçersiz e-posta formatı.")
	}

	emailExists, _ := s.repo.ExistsByEmail(email)
	if emailExists {
		return errors.New("Bu e-posta adresi zaten kullanımda.")
	}

	hashedPassword, _ := utils.HashPassword(password)
	user := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("Girdiğiniz e-posta adresine ait bir hesap bulunamadı.")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("Hatalı şifre girdiniz. Lütfen tekrar deneyin.")
	}

	return utils.GenerateToken(uint(user.ID), s.jwtSecret)
}

func (s *UserService) Update(userID uint, username, email string) error {
	if !utils.IsValidEmail(email) {
		return errors.New("Geçersiz e-posta formatı.")
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Username = username
	user.Email = email
	return s.repo.Update(user)
}

func (s *UserService) Delete(userID uint) error {
	return s.repo.Delete(userID)
}
