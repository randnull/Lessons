package gRPC_client

import (
	"context"
	"github.com/randnull/Lessons/internal/models"
)

type GRPCClientInt interface {
	//GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.NewUser) (string, error)
	//GetAllUsers(ctx context.Context) (*pb.GetAllResponse, error)
	Close()
}
