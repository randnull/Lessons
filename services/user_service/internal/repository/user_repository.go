package repository

import "github.com/randnull/Lessons/internal/models"

type UserRepository interface {
	CreateUser(user *models.CreateUser) (string, error)
	GetUserInfoById(telegramID int64) (*models.UserDB, error)
	//CheckExistUser(user_id string) bool
}
