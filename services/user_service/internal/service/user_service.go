package service

import (
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
)

type UserServiceInt interface {
	GetUserById(UserId string) (*models.UserDB, error)
	GetStudentById(UserId string) (*models.UserDB, error)
	GetTutorById(TutorID string) (*models.TutorDB, error)
	//GetUserByTelegramId(TelegramId int64) (*models.UserDB, error)
	UpdateNameTutor(tutorID string, name string) error
	CreateUser(user models.CreateUser) (string, error)
	GetTutors() ([]*pb.Tutor, error)
	GetTutorsPagination(page int, size int) (*pb.GetTutorsPaginationResponse, error)
	UpdateBioTutor(userID string, bio string) error
	UpdateTutorTags(tutorID string, tags []string) error
	CreateReview(tutorID, studentID string, rating int, comment string) (string, error)
	GetReviews(tutorID string) ([]models.Review, error)
	GetReviewById(reviewID string) (*models.Review, error)
	GetTutorInfoById(tutorID string) (*models.TutorDetails, error)
	ChangeTutorActive(tutorID string, IsActive bool) error
	CreateNewResponse(tutorID string) error
	AddResponses(tutorID int64, responseCount int) (int, error)
}

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserServiceInt {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s *UserService) GetUserById(UserId string) (*models.UserDB, error) {
	return s.userRepository.GetUserById(UserId)
}

func (s *UserService) GetStudentById(UserId string) (*models.UserDB, error) {
	return s.userRepository.GetStudentById(UserId)
}

//
//func (s *UserService) GetUserByTelegramId(TelegramId int64) (*models.UserDB, error) {
//	return s.userRepository.GetUserByTelegramId(TelegramId)
//}

func (s *UserService) CreateUser(user models.CreateUser) (string, error) {
	return s.userRepository.CreateUser(&user)
}

func (s *UserService) GetTutors() ([]*pb.Tutor, error) {
	return s.userRepository.GetAllTutors()
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

func (s *UserService) GetTutorsPagination(page int, size int) (*pb.GetTutorsPaginationResponse, error) {
	limit := size
	offset := (page - 1) * size

	tutors, count, err := s.userRepository.GetAllTutorsPagination(limit, offset)

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
func (s *UserService) CreateReview(tutorID string, studentID string, rating int, comment string) (string, error) {
	return s.userRepository.CreateReview(tutorID, studentID, rating, comment)
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

	tags, err := s.userRepository.GetTagsByTutorID(tutorID)
	if err != nil {
		return nil, err
	}

	return &models.TutorDetails{
		Tutor:   *tutor,
		Reviews: reviews,
		Tags:    tags,
	}, nil
}

func (s *UserService) ChangeTutorActive(tutorID string, IsActive bool) error {
	return s.userRepository.SetNewIsActiveTutor(tutorID, IsActive)
}

func (s *UserService) CreateNewResponse(tutorID string) error {
	return s.userRepository.RemoveOneResponse(tutorID)
}

func (s *UserService) AddResponses(tutorID int64, responseCount int) (int, error) {
	return s.userRepository.AddResponses(tutorID, responseCount)
}
