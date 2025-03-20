package gRPC_client

import (
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
)

type GRPCClientInt interface {
	GetUser(userID string) (*models.User, error)
	GetStudent(userID string) (*models.User, error)
	GetUserByTelegramID(telegramID int64) (*models.User, error)
	//CreateUser(ctx context.Context, user *models.CreateUser) (string, error)
	GetAllUsers() (*pb.GetAllResponse, error)
	Close()
}
