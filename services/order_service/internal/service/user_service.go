package service

import (
	"context"
	"errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
)

type UserServiceInt interface {
	GetUser(UserID string) (*models.User, error)
	GetTutor(TutorID string) (*models.Tutor, error)

	GetAllTutorsPagination(page int, size int, tag string) (*models.TutorsPagination, error)
	UpdateBioTutor(BioModel models.UpdateBioTutor, UserData models.UserData) error
	UpdateTagsTutor(tags []string, UserData models.UserData) (bool, error)
	CreateReview(orderID string, tutorID string, comment string, rating int, UserData models.UserData) (string, error)
	GetReviewsByTutor(tutorID string) ([]models.Review, error)
	GetReviewsByID(reviewID string) (*models.Review, error)
	GetTutorInfoById(tutorID string) (*models.TutorDetails, error)
	ChangeTutorActive(isActive bool, UserData models.UserData) (bool, error)
	UpdateTutorName(tutorID, name string) error
}

type UserService struct {
	GRPCClient     gRPC_client.GRPCClientInt
	ProducerBroker rabbitmq.RabbitMQInterface
}

func NewUSerService(grpcClient gRPC_client.GRPCClientInt, producerBroker rabbitmq.RabbitMQInterface) UserServiceInt {
	return &UserService{
		GRPCClient:     grpcClient,
		ProducerBroker: producerBroker,
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
		logger.Error("[UserService] GetAllUsers error GetAllUsers: " + err.Error())
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
	_, err := u.GRPCClient.UpdateBioTutor(context.Background(), BioModel.Bio, UserData.UserID)
	//success
	if err != nil {

		return err
	}

	return nil
}

func (u *UserService) GetAllTutorsPagination(page int, size int, tag string) (*models.TutorsPagination, error) {
	usersRPC, err := u.GRPCClient.GetTutorsPagination(context.Background(), page, size, tag)

	if err != nil {
		logger.Error("[UserService] GetAllTutorsPagination error GetTutorsPagination: " + err.Error())
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

func (u *UserService) UpdateTagsTutor(tags []string, UserData models.UserData) (bool, error) {
	success, err := u.GRPCClient.UpdateTagsTutor(context.Background(), tags, UserData.UserID)
	if err != nil {
		return false, err
	}

	if success {
		ChangeTagsToBroker := &models.ChangeTagsTutorToBroker{
			TutorTelegramID: UserData.TelegramID,
			Tags:            tags,
		}

		err := u.ProducerBroker.Publish("tutors_tags_change", ChangeTagsToBroker)
		
		if err != nil {
			logger.Error("[OrderService] ResponseToOrder Error publishing order: " + err.Error())
		}
	}

	return success, nil
}

func (u *UserService) CreateReview(orderID string, tutorID string, comment string, rating int, UserData models.UserData) (string, error) {
	// вот тут обязательные провекрки
	return u.GRPCClient.CreateReview(context.Background(), orderID, tutorID, comment, rating)
}

func (u *UserService) GetReviewsByTutor(tutorID string) ([]models.Review, error) {
	reviews, err := u.GRPCClient.GetReviewsByTutor(context.Background(), tutorID)
	if err != nil {
		logger.Error("[UserService] GetReviewsByTutor error GetReviewsByTutor: " + err.Error())
		return nil, err
	}

	return reviews, nil
}

func (u *UserService) GetReviewsByID(reviewID string) (*models.Review, error) {
	review, err := u.GRPCClient.GetReviewsByID(context.Background(), reviewID)
	if err != nil {
		logger.Error("[UserService] GetReviewsByID error GetReviewsByID: " + err.Error())

		return nil, err
	}

	return review, nil
}

func (u *UserService) GetTutorInfoById(tutorID string) (*models.TutorDetails, error) {
	TutorDetails, err := u.GRPCClient.GetTutorInfoById(context.Background(), tutorID)

	if err != nil {
		logger.Error("[UserService] GetTutorInfoById error GetTutorInfoById: " + err.Error())
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
