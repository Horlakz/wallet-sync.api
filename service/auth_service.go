package service

import (
	"errors"

	"github.com/horlakz/wallet-sync.api/dto"
	"github.com/horlakz/wallet-sync.api/internal/helper"
	"github.com/horlakz/wallet-sync.api/model"
	user_repository "github.com/horlakz/wallet-sync.api/repository/user"
)

type AuthServiceInterface interface {
	Login(email, password string) (string, error)
	Register(data dto.RegisterDTO) error
}

type authService struct {
	userRepo user_repository.UserRepository
	encrpyt  helper.HashingInterface
	jwt      helper.JwtInterface
}

func NewAuthService(userRepo user_repository.UserRepository) AuthServiceInterface {
	return &authService{
		userRepo: userRepo,
		encrpyt:  helper.NewHashing(),
		jwt:      helper.NewJwt(),
	}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	match, err := s.encrpyt.ComparePassword(password, user.Password)
	if err != nil {
		return "", err
	}

	if !match {
		return "", errors.New("invalid credentials")
	}

	accessToken, err := s.jwt.CreateToken(user.ID.String(), "access")
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *authService) Register(data dto.RegisterDTO) error {
	hashedPassword, err := s.encrpyt.HashPassword(data.Password)
	if err != nil {
		return err
	}
	user := &model.User{
		Email:    data.Email,
		Name:     data.Name,
		Password: hashedPassword,
	}

	if existingUser, _ := s.userRepo.FindByEmail(data.Email); existingUser != nil {
		return errors.New("user already exists")
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	return nil
}
