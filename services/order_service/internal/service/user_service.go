package service

import (
	"context"
	"errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
)

type UserServiceInt interface {
	GetUser(UserID string) (*models.User, error)
	GetTutor(TutorID string) (*models.Tutor, error)

	GetAllTutorsPagination(page int, size int) (*models.TutorsPagination, error)
	UpdateBioTutor(BioModel models.UpdateBioTutor, UserData models.UserData) error
	UpdateTagsTutor(tags []string, TutorID string) (bool, error)
	CreateReview(studentID string, tutorID string, comment string, rating int) (string, error)
	GetReviewsByTutor(tutorID string) ([]models.Review, error)
	GetReviewsByID(reviewID string) (*models.Review, error)
	GetTutorInfoById(tutorID string) (*models.TutorDetails, error)
	ChangeTutorActive(isActive bool, UserData models.UserData) (bool, error)
	UpdateTutorName(tutorID, name string) error
}

type UserService struct {
	GRPCClient gRPC_client.GRPCClientInt
}

func NewUSerService(grpcClient gRPC_client.GRPCClientInt) UserServiceInt {
	return &UserService{
		GRPCClient: grpcClient,
	}
}

func (u *UserService) GetTutor(TutorID string) (*models.Tutor, error) {
	return u.GRPCClient.GetTutor(context.Background(), TutorID)
}

func (u *UserService) GetUser(UserID string) (*models.User, error) {
	return u.GRPCClient.GetUser(context.Background(), UserID)
}

func (u *UserService) GetAllUsers() ([]*models.TutorForList, error) {
	usersRPC, err := u.GRPCClient.GetAllUsers(context.Background())

	if err != nil {
		return nil, err
	}

	var users []*models.TutorForList

	for _, grpcUser := range usersRPC.Tutors {
		users = append(users, &models.TutorForList{
			Id:   grpcUser.User.Id,
			Name: grpcUser.User.Name,
			Tags: grpcUser.Tags,
		})
	}

	return users, nil
}

func (u *UserService) UpdateBioTutor(BioModel models.UpdateBioTutor, UserData models.UserData) error {
	success, err := u.GRPCClient.UpdateBioTutor(context.Background(), BioModel.Bio, UserData.UserID)

	if !success {
		return err
	}

	return nil
}

func (u *UserService) GetAllTutorsPagination(page int, size int) (*models.TutorsPagination, error) {
	usersRPC, err := u.GRPCClient.GetTutorsPagination(context.Background(), page, size)

	if err != nil {
		return nil, err
	}

	var tutors []*models.TutorForList

	for _, grpcUser := range usersRPC.Tutors {
		tutors = append(tutors, &models.TutorForList{
			Id:   grpcUser.User.Id,
			Name: grpcUser.User.Name,
			Tags: grpcUser.Tags,
		})
	}

	addPage := 0

	if int(usersRPC.Count)%size != 0 {
		addPage += 1
	}

	return &models.TutorsPagination{
		Tutors: tutors,
		Pages:  (int(usersRPC.Count) / size) + addPage,
	}, nil
}

func (u *UserService) UpdateTagsTutor(tags []string, TutorID string) (bool, error) {
	success, err := u.GRPCClient.UpdateTagsTutor(context.Background(), tags, TutorID)
	if err != nil {
		return false, err
	}
	return success, nil
}

func (u *UserService) CreateReview(studentID string, tutorID string, comment string, rating int) (string, error) {
	return u.GRPCClient.CreateReview(context.Background(), studentID, tutorID, comment, rating)
}

func (u *UserService) GetReviewsByTutor(tutorID string) ([]models.Review, error) {
	reviews, err := u.GRPCClient.GetReviewsByTutor(context.Background(), tutorID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (u *UserService) GetReviewsByID(reviewID string) (*models.Review, error) {
	review, err := u.GRPCClient.GetReviewsByID(context.Background(), reviewID)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (u *UserService) GetTutorInfoById(tutorID string) (*models.TutorDetails, error) {
	TutorDetails, err := u.GRPCClient.GetTutorInfoById(context.Background(), tutorID)

	if err != nil {
		return nil, err
	}

	return TutorDetails, nil
}

func (u *UserService) ChangeTutorActive(isActive bool, UserData models.UserData) (bool, error) {
	isOk, err := u.GRPCClient.ChangeTutorActive(context.Background(), UserData.UserID, isActive)

	if err != nil {
		return false, err
	}

	return isOk, nil
}

func (u *UserService) UpdateTutorName(tutorID, name string) error {
	isOk, err := u.GRPCClient.UpdateNameTutor(context.Background(), tutorID, name)

	if err != nil || !isOk {
		return errors.New("cannot update tutor name")
	}

	return nil
}
