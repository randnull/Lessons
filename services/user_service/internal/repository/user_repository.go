package repository

import (
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.CreateUser) (string, error)
	GetUserById(userID string) (*models.UserDB, error)
	GetUserByTelegramId(telegramID int64) (*models.UserDB, error)
	GetAllUsers() ([]*pb.User, error)
	//CheckExistUser(user_id string) bool
}
