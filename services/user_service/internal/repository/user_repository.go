package repository

import (
	"github.com/randnull/Lessons/internal/models"
	pb "github.com/randnull/Lessons/pkg/gRPC"
)

type UserRepository interface {
	CreateUser(user *models.CreateUser) (string, error)

	GetStudentById(userID string) (*models.UserDB, error)
	GetTutorByID(userID string) (*models.TutorDB, error)
	GetUserById(userID string) (*models.UserDB, error)
	GetUserByTelegramId(telegramID int64, userRole string) (*models.UserDB, error)
	GetAllUsers() ([]*pb.User, error)
	GetAllTutorsResponseCondition(minResponseCount int) ([]*models.TutorWithResponse, error)
	GetAllTutorsPagination(limit int, offset int, tag string) ([]*pb.Tutor, int, error)

	UpdateTutorBio(userID string, bio string) error
	UpdateTutorTags(tutorID string, tags []string) error
	UpdateTutorName(tutorID string, name string) error

	CreateReview(tutorID, orderID string, rating int, comment string) (string, error)
	GetReviews(tutorID string) ([]models.Review, error)
	GetReviewById(reviewID string) (*models.Review, error)
	SetReviewActive(reviewID string, tutorID string) error

	SetNewIsActiveTutor(tutorID string, IsActive bool) error

	AddResponses(tutorTelegramID int64, responseCount int) (int, error)
	RemoveOneResponse(tutorID string) error

	BanUser(telegramID int64, isBanned bool) error
}
