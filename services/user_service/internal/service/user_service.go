package service

import (
	"fmt"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/pkg/custom_errors"
	pb "github.com/randnull/Lessons/pkg/gRPC"
	lg "github.com/randnull/Lessons/pkg/logger"
)

type UserServiceInt interface {
	CreateUser(user models.CreateUser) (string, error)

	GetUserById(UserId string) (*models.UserDB, error)
	GetStudentById(UserId string) (*models.UserDB, error)
	GetTutorById(TutorID string) (*models.TutorDB, error)
	GetUserByTelegramId(TelegramId int64, userRole string) (*models.UserDB, error)
	GetAllUsers() ([]*pb.User, error)
	GetTutorsPagination(page int, size int, tag string) (*pb.GetTutorsPaginationResponse, error)
	GetTutorInfoById(tutorID string) (*models.TutorDetails, error)

	GetReviews(tutorID string) ([]models.Review, error)
	GetReviewById(reviewID string) (*models.Review, error)
	CreateReview(tutorID, orderID string, rating int, comment string) (string, error)
	SetReviewActive(reviewID string) error

	CreateNewResponse(tutorID string) error
	AddResponses(tutorID int64, responseCount int) (int, error)

	ChangeTutorActive(tutorID string, IsActive bool) error
	UpdateNameTutor(tutorID string, name string) error
	UpdateBioTutor(userID string, bio string) error
	UpdateTutorTags(tutorID string, tags []string) error

	BanUser(telegramID int64, isBanned bool) error
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceInt {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s *UserService) GetStudentById(UserId string) (*models.UserDB, error) {
	return s.userRepository.GetStudentById(UserId)
}

func (s *UserService) GetUserByTelegramId(TelegramId int64, userRole string) (*models.UserDB, error) {
	return s.userRepository.GetUserByTelegramId(TelegramId, userRole)
}

func (s *UserService) CreateUser(user models.CreateUser) (string, error) {
	if (user.Role != models.RoleStudent) && (user.Role != models.RoleTutor) && (user.Role != models.RoleAdmin) {
		return "", custom_errors.ErrorIncorrectRole
	}

	return s.userRepository.CreateUser(&user)
}

func (s *UserService) UpdateBioTutor(userID string, bio string) error {
	return s.userRepository.UpdateTutorBio(userID, bio)
}

func (s *UserService) UpdateNameTutor(tutorID string, name string) error {
	return s.userRepository.UpdateTutorName(tutorID, name)
}

func (s *UserService) GetTutorById(TutorID string) (*models.TutorDB, error) {
	return s.userRepository.GetTutorByID(TutorID)
}

func (s *UserService) GetTutorsPagination(page int, size int, tag string) (*pb.GetTutorsPaginationResponse, error) {
	limit := size
	offset := (page - 1) * size

	tutors, count, err := s.userRepository.GetAllTutorsPagination(limit, offset, tag)

	if err != nil {
		return nil, err
	}

	return &pb.GetTutorsPaginationResponse{
		Count:  int64(count),
		Tutors: tutors,
	}, nil
}

func (s *UserService) UpdateTutorTags(tutorID string, tags []string) error {
	return s.userRepository.UpdateTutorTags(tutorID, tags)
}

func (s *UserService) CreateReview(tutorID string, orderID string, rating int, comment string) (string, error) {
	return s.userRepository.CreateReview(tutorID, orderID, rating, comment)
}
func (s *UserService) GetReviews(tutorID string) ([]models.Review, error) {
	return s.userRepository.GetReviews(tutorID)
}
func (s *UserService) GetReviewById(reviewID string) (*models.Review, error) {
	return s.userRepository.GetReviewById(reviewID)
}

func (s *UserService) GetTutorInfoById(tutorID string) (*models.TutorDetails, error) {
	tutor, err := s.userRepository.GetTutorByID(tutorID)
	if err != nil {
		return nil, err
	}

	reviews, err := s.userRepository.GetReviews(tutorID)
	if err != nil {
		return nil, err
	}

	return &models.TutorDetails{
		Tutor:   *tutor,
		Reviews: reviews,
	}, nil
}

func (s *UserService) ChangeTutorActive(tutorID string, IsActive bool) error {
	return s.userRepository.SetNewIsActiveTutor(tutorID, IsActive)
}

func (s *UserService) CreateNewResponse(tutorID string) error {
	return s.userRepository.RemoveOneResponse(tutorID)
}

func (s *UserService) AddResponses(tutorID int64, responseCount int) (int, error) {
	if responseCount < 1 {
		return 0, custom_errors.ErrorCountLessZero
	}

	return s.userRepository.AddResponses(tutorID, responseCount)
}

func (s *UserService) SetReviewActive(reviewID string) error {
	review, err := s.userRepository.GetReviewById(reviewID)

	if err != nil {
		lg.Error(fmt.Sprintf("error with get review for reviewID: %v. Error: %v", reviewID, err.Error()))
		return err
	}

	return s.userRepository.SetReviewActive(reviewID, review.TutorID)
}

func (s *UserService) BanUser(telegramID int64, isBanned bool) error {
	return s.userRepository.BanUser(telegramID, isBanned)
}

func (s *UserService) GetUserById(UserId string) (*models.UserDB, error) {
	return s.userRepository.GetUserById(UserId)
}

func (s *UserService) GetAllUsers() ([]*pb.User, error) {
	return s.userRepository.GetAllUsers()
}
