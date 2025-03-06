package gRPC_client

import (
	"context"
	"github.com/randnull/Lessons/internal/models"
)

type GRPCClientInt interface {
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.CreateUser) (string, error)
	Close()
}
