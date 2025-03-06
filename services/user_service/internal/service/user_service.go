package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
)

type UserServiceInt interface {
	GetUserById(telegramID int64) (*models.UserDB, error)
	CreateUser(user models.CreateUser) (string, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceInt {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s UserService) GetUserById(telegramID int64) (*models.UserDB, error) {
	return s.userRepository.GetUserInfoById(telegramID)
}

func (s UserService) CreateUser(user models.CreateUser) (string, error) {
	return s.userRepository.CreateUser(&user)
}
