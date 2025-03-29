package gRPC_client

import (
	"context"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
)

type GRPCClientInt interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetStudent(ctx context.Context, userID string) (*models.User, error)
	GetTutor(ctx context.Context, TutorID string) (*models.Tutor, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	GetAllUsers(ctx context.Context) (*pb.GetAllResponse, error)
	GetTutorsPagination(ctx context.Context, page int, size int) (*pb.GetAllResponse, error)
	UpdateBioTutor(ctx context.Context, bio string, TutorID string) (bool, error)
	Close()
}
