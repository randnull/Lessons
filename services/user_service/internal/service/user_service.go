package service

import (
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
)

type UserServiceInt interface {
	GetUserById(userID string) (*models.User, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceInt {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s UserService) GetUserById(userID string) (*models.User, error) {
	return &models.User{
		UserId: "c1f6d20d-35aa-4268-b711-fd19f994c0ae",
		Name:   "John",
	}, nil
}
