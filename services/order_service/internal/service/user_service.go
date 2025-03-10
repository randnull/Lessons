package service

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type UserServiceInt interface {
	CreateUser(InitData initdata.InitData) (string, error)
	GetUser(TelegramID int64) (*models.User, error)
}

type UserService struct {
	GRPCClient gRPC_client.GRPCClientInt
}

func NewUSerService(grpcClient gRPC_client.GRPCClientInt) UserServiceInt {
	return &UserService{
		GRPCClient: grpcClient,
	}
}

func (u *UserService) GetUser(TelegramID int64) (*models.User, error) {
	return u.GRPCClient.GetUser(context.Background(), TelegramID)
}

func (u *UserService) CreateUser(InitData initdata.InitData) (string, error) {
	fmt.Println(InitData)
	NewUser := &models.CreateUser{
		Name:       InitData.User.FirstName,
		TelegramId: InitData.User.ID,
	}

	userID, err := u.GRPCClient.CreateUser(context.Background(), NewUser)
	if err != nil {
		return "", err
	}

	return userID, nil
}
