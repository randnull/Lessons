package gRPC_client

import (
	"context"
	"github.com/randnull/Lessons/internal/models"
	pb "github.com/randnull/Lessons/pkg/gRPC"
)

type GRPCClientInt interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetStudent(ctx context.Context, userID string) (*models.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	GetAllUsers(ctx context.Context) (*pb.GetAllResponse, error)

	GetTutor(ctx context.Context, tutorID string) (*models.Tutor, error)
	GetTutorInfoById(ctx context.Context, tutorID string, isOwn bool) (*models.TutorDetails, error)
	GetTutorsPagination(ctx context.Context, page, size int, tag string) (*pb.GetTutorsPaginationResponse, error)
	UpdateBioTutor(ctx context.Context, bio, tutorID string) (bool, error)
	UpdateTagsTutor(ctx context.Context, tags []string, tutorID string) (bool, error)
	UpdateNameTutor(ctx context.Context, tutorID, name string) (bool, error)

	ChangeTutorActive(ctx context.Context, tutorID string, active bool) (bool, error)
	CreateNewResponse(ctx context.Context, tutorID string) (bool, error)

	CreateReview(ctx context.Context, orderID, tutorID, comment string, rating int) (string, error)
	GetReviewsByTutor(ctx context.Context, tutorID string) ([]models.Review, error)
	GetReviewsByID(ctx context.Context, reviewID string) (*models.Review, error)
	SetActiveToReview(ctx context.Context, reviewID string) (bool, error)

	BanUser(ctx context.Context, telegramID int64, isBanned bool) (bool, error)

	Close() error
}
