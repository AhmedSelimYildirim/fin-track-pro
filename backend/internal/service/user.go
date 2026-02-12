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
		return errors.New("Gecersiz e-posta formati.")
	}
	emailExists, _ := s.repo.ExistsByEmail(email)
	if emailExists {
		return errors.New("Bu e-posta adresi zaten kullanimda.")
	}
	hashedPassword, _ := utils.HashPassword(password)
	user := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("Girdiginiz e-posta adresine ait bir hesap bulunamadi.")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("Hatali sifre girdiniz. Lutfen tekrar deneyin.")
	}
	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	return token, user, err
}

func (s *UserService) Update(userID int64, username, email, password string) error {
	if !utils.IsValidEmail(email) {
		return errors.New("Gecersiz e-posta formati.")
	}
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.Email != email {
		exists, _ := s.repo.ExistsByEmail(email)
		if exists {
			return errors.New("Bu e-posta adresi zaten kullanimda.")
		}
	}
	user.Username = username
	user.Email = email
	if password != "" {
		hashedPassword, _ := utils.HashPassword(password)
		user.Password = hashedPassword
	}
	return s.repo.Update(user)
}

func (s *UserService) Delete(userID int64) error {
	return s.repo.Delete(userID)
}
