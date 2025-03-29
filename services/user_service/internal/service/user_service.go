package service

import (
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
	"log"
)

type UserServiceInt interface {
	GetUserById(UserId string) (*models.UserDB, error)
	GetStudentById(UserId string) (*models.UserDB, error)
	GetTutorById(TutorID string) (*models.TutorDB, error)
	//GetUserByTelegramId(TelegramId int64) (*models.UserDB, error)
	CreateUser(user models.CreateUser) (string, error)
	GetTutors() ([]*pb.User, error)
	GetTutorsPagination(page int, size int) ([]*pb.User, error)
	UpdateBioTutor(userID string, bio string) error
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

func (s *UserService) GetStudentById(UserId string) (*models.UserDB, error) {
	return s.userRepository.GetStudentById(UserId)
}

//
//func (s *UserService) GetUserByTelegramId(TelegramId int64) (*models.UserDB, error) {
//	return s.userRepository.GetUserByTelegramId(TelegramId)
//}

func (s *UserService) CreateUser(user models.CreateUser) (string, error) {
	log.Println("USE!!", user)
	return s.userRepository.CreateUser(&user)
}

func (s *UserService) GetTutors() ([]*pb.User, error) {
	return s.userRepository.GetAllTutors()
}

func (s *UserService) UpdateBioTutor(userID string, bio string) error {
	return s.userRepository.UpdateTutorBio(userID, bio)
}

func (s *UserService) GetTutorById(TutorID string) (*models.TutorDB, error) {
	return s.userRepository.GetTutorByID(TutorID)
}

func (s *UserService) GetTutorsPagination(page int, size int) ([]*pb.User, error) {
	limit := size
	offset := (page - 1) * size

	return s.userRepository.GetAllTutorsPagination(limit, offset)
}
