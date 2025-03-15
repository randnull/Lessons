package service

import (
	"fmt"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
)

type UserServiceInt interface {
	GetUserById(UserId string) (*models.UserDB, error)
	//GetUserByTelegramId(TelegramId int64) (*models.UserDB, error)
	CreateUser(user models.CreateUser) (string, error)
	GetAllUsers() ([]*pb.User, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceInt {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s *UserService) GetUserById(UserId string) (*models.UserDB, error) {
	return s.userRepository.GetUserById(UserId)
}

//
//func (s *UserService) GetUserByTelegramId(TelegramId int64) (*models.UserDB, error) {
//	return s.userRepository.GetUserByTelegramId(TelegramId)
//}

func (s *UserService) CreateUser(user models.CreateUser) (string, error) {
	fmt.Println("USE!!", user)
	return s.userRepository.CreateUser(&user)
}

func (s *UserService) GetAllUsers() ([]*pb.User, error) {
	return s.userRepository.GetAllUsers()
}
