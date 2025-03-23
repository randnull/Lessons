package repository

import (
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.CreateUser) (string, error)
	GetStudentById(userID string) (*models.UserDB, error)
	GetUserById(userID string) (*models.UserDB, error)
	GetTutorByID(userID string) (*models.TutorDB, error)
	GetUserByTelegramId(telegramID int64, userRole string) (*models.UserDB, error)
	GetAllUsers() ([]*pb.User, error)
	UpdateTutorBio(userID string, bio string) error
	//CheckExistUser(user_id string) bool
}
