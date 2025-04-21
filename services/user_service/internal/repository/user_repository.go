package repository

import (
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.CreateUser) (string, error)
	GetStudentById(userID string) (*models.UserDB, error)
	GetUserById(userID string) (*models.UserDB, error)
	GetTutorByID(userID string) (*models.TutorDB, error)
	GetUserByTelegramId(telegramID int64, userRole string) (*models.UserDB, error)
	GetAllTutors() ([]*pb.Tutor, error)
	GetAllTutorsPagination(limit int, offset int, tag string) ([]*pb.Tutor, int, error)
	UpdateTutorBio(userID string, bio string) error
	UpdateTutorTags(tutorID string, tags []string) error
	UpdateTutorName(tutorID string, name string) error
	CreateReview(tutorID, orderID string, rating int, comment string) (string, error)
	GetReviews(tutorID string) ([]models.Review, error)
	GetReviewById(reviewID string) (*models.Review, error)
	GetTagsByTutorID(tutorID string) ([]string, error)
	SetNewIsActiveTutor(tutorID string, IsActive bool) error
	AddResponses(tutorTelegramID int64, responseCount int) (int, error)
	RemoveOneResponse(tutorID string) error
}
