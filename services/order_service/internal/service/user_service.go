package service

import (
	"context"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
)

type UserServiceInt interface {
	//CreateUser(UserData models.UserData) (string, error)
	GetUser(UserID string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
}

type UserService struct {
	GRPCClient gRPC_client.GRPCClientInt
}

func NewUSerService(grpcClient gRPC_client.GRPCClientInt) UserServiceInt {
	return &UserService{
		GRPCClient: grpcClient,
	}
}

func (u *UserService) GetUser(UserID string) (*models.User, error) {
	return u.GRPCClient.GetUser(context.Background(), UserID)
}

//func (u *UserService) CreateUser(UserData models.UserData) (string, error) {
//	NewUser := &models.CreateUser{
//		Name:       UserData.FirstName,
//		TelegramId: UserData.TelegramID,
//	}
//
//	userID, err := u.GRPCClient.CreateUser(context.Background(), NewUser)
//	if err != nil {
//		return "", err
//	}
//
//	return userID, nil
//}

func (u *UserService) GetAllUsers() ([]*models.User, error) {
	usersRPC, err := u.GRPCClient.GetAllUsers(context.Background())

	if err != nil {
		return nil, err
	}

	var users []*models.User

	for _, grpcUser := range usersRPC.Users {
		users = append(users, &models.User{
			Id:   grpcUser.Id,
			Name: grpcUser.Name,
		})
	}

	return users, nil
}
