package repository

import "github.com/randnull/Lessons/internal/models"

type UserRepository interface {
	//CreateUser(user *models.User) (string, error)
	GetUserById(user_id string) (*models.User, error)
	//CheckExistUser(user_id string) bool
}
