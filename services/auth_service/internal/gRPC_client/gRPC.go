package gRPC_client

import (
	"context"
	"github.com/randnull/Lessons/internal/models"
)

type GRPCClientInt interface {
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.NewUser) (string, error)
	Close()
}
