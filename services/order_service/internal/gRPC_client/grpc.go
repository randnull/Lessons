package gRPC_client

import (
	"context"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
)

type GRPCClientInt interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	//CreateUser(ctx context.Context, user *models.CreateUser) (string, error)
	GetAllUsers(ctx context.Context) (*pb.GetAllResponse, error)
	Close()
}
