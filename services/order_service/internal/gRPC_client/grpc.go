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
	GetTutorsPagination(ctx context.Context, page int, size int) (*pb.GetTutorsPaginationResponse, error)
	UpdateBioTutor(ctx context.Context, bio string, TutorID string) (bool, error)
	UpdateTagsTutor(ctx context.Context, tags []string, TutorID string) (bool, error)
	CreateReview(ctx context.Context, studentID string, tutorID string, comment string, rating int) (string, error)
	GetReviewsByTutor(ctx context.Context, tutorID string) ([]models.Review, error)
	GetReviewsByID(ctx context.Context, reviewID string) (*models.Review, error)
	GetTutorInfoById(ctx context.Context, TutorID string) (*models.TutorDetails, error)
	Close()
}
