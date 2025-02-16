package gRPC_client

import (
	"context"
	pb "github.com/randnull/Lessons/internal/gRPC"
)

type GRPCClientInt interface {
	GetUser(ctx context.Context, userID string) (*pb.User, error)
	Close()
}
