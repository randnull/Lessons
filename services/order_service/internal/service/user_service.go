package service

import (
	"context"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
)

type UserServiceInt interface {
	GetUser(UserID string) (*models.User, error)
	GetTutor(TutorID string) (*models.Tutor, error)

	GetAllUsers() ([]*models.User, error)
	UpdateBioTutor(BioModel models.UpdateBioTutor, UserData models.UserData) error
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

func (u *UserService) GetAllUsers() ([]*models.User, error) {
	usersRPC, err := u.GRPCClient.GetAllUsers(context.Background())

	if err != nil {
		return nil, err
	}

	var users []*models.User

	for _, grpcUser := range usersRPC.Users {
		users = append(users, &models.User{
			Id:   grpcUser.Id,
			Name: grpcUser.Name,
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
