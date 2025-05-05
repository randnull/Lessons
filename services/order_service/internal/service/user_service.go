package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/gRPC_client"
	lg "github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/utils"
	"strconv"
	"time"
)

type UserServiceInt interface {
	GetTutor(TutorID string, UserData models.UserData) (*models.Tutor, error)
	GetAllTutorsPagination(page int, size int, tag string, UserData models.UserData) (*models.TutorsPagination, error)
	UpdateBioTutor(BioModel models.UpdateBioTutor, UserData models.UserData) error
	UpdateTagsTutor(tags []string, UserData models.UserData) (bool, error)
	CreateReview(ReviewRequest models.ReviewRequest, UserData models.UserData) (string, error)
	GetReviewsByTutor(tutorID string, UserData models.UserData) ([]models.Review, error)
	GetReviewsByID(reviewID string, UserData models.UserData) (*models.Review, error)
	GetTutorInfoById(tutorID string, UserData models.UserData) (*models.TutorDetails, error)
	ChangeTutorActive(isActive bool, UserData models.UserData) (bool, error)
	UpdateTutorName(tutorID, name string, UserData models.UserData) error
	GetAllUsers(UserData models.UserData) ([]*models.User, error)
	GetUserById(userID string, UserData models.UserData) (*models.User, error)
	SetReviewActive(reviewID string, UserData models.UserData) error
	BanUser(banUser models.BanUser, UserData models.UserData) error
}

type UserService struct {
	GRPCClient      gRPC_client.GRPCClientInt
	ProducerBroker  rabbitmq.RabbitMQInterface
	orderRepository repository.OrderRepository
}

func NewUserService(grpcClient gRPC_client.GRPCClientInt, producerBroker rabbitmq.RabbitMQInterface, orderRepo repository.OrderRepository) UserServiceInt {
	return &UserService{
		GRPCClient:      grpcClient,
		ProducerBroker:  producerBroker,
		orderRepository: orderRepo,
	}
}

func (u *UserService) GetTutor(TutorID string, UserData models.UserData) (*models.Tutor, error) {
	return u.GRPCClient.GetTutor(context.Background(), TutorID)
}

func (u *UserService) GetAllUsers(UserData models.UserData) ([]*models.User, error) {
	usersRPC, err := u.GRPCClient.GetAllUsers(context.Background())

	if err != nil {
		lg.Error("[UserService] GetAllUsers error GetAllUsers: " + err.Error())
		return nil, err
	}

	var users []*models.User

	for _, grpcUser := range usersRPC.Users {
		users = append(users, &models.User{
			Id:         grpcUser.Id,
			Name:       grpcUser.Name,
			TelegramID: grpcUser.TelegramId,
			Role:       grpcUser.Role,
		})
	}

	return users, nil
}

func (u *UserService) UpdateBioTutor(BioModel models.UpdateBioTutor, UserData models.UserData) error {
	if utils.ContainsBadWords(BioModel.Bio) {
		return custom_errors.ErrorBanWords
	}

	_, err := u.GRPCClient.UpdateBioTutor(context.Background(), BioModel.Bio, UserData.UserID)
	//success
	if err != nil {

		return err
	}

	return nil
}

func (u *UserService) GetAllTutorsPagination(page int, size int, tag string, UserData models.UserData) (*models.TutorsPagination, error) {
	usersRPC, err := u.GRPCClient.GetTutorsPagination(context.Background(), page, size, tag)

	if err != nil {
		lg.Error("[UserService] GetAllTutorsPagination error GetTutorsPagination: " + err.Error())
		return nil, err
	}

	var tutors []*models.TutorForList

	for _, grpcUser := range usersRPC.Tutors {
		tutors = append(tutors, &models.TutorForList{
			Id:     grpcUser.User.Id,
			Name:   grpcUser.User.Name,
			Tags:   grpcUser.Tags,
			Rating: grpcUser.Rating,
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

	for _, tag := range tags {
		if utils.ContainsBadWords(tag) {
			return false, custom_errors.ErrorBanWords
		}
	}

	if err != nil {
		lg.Error("[UserService] GetAllTutorsPagination error UpdateTagsTutor: " + err.Error())
		return false, err
	}

	if success {
		ChangeTagsToBroker := &models.ChangeTagsTutorToBroker{
			TutorTelegramID: UserData.TelegramID,
			Tags:            tags,
		}

		err := u.ProducerBroker.Publish("tutors_tags_change", ChangeTagsToBroker)

		if err != nil {
			lg.Error("[UserService] ResponseToOrder Error publishing order: " + err.Error())
		}
	}

	return success, nil
}

func (u *UserService) CreateReview(ReviewRequest models.ReviewRequest, UserData models.UserData) (string, error) {
	if utils.ContainsBadWords(ReviewRequest.Comment) {
		return "", custom_errors.ErrorBanWords
	}

	response, err := u.orderRepository.GetResponseById(ReviewRequest.ResponseID)

	if err != nil {
		lg.Error("[UserService] CreateReview error GetResponseById: " + err.Error())
		return "", custom_errors.ErrorGetResponse
	}

	order, err := u.orderRepository.GetOrderByID(response.OrderID)

	if err != nil {
		lg.Error("[UserService] CreateReview error GetOrderByID: " + err.Error())
		return "", custom_errors.ErrorServiceError
	}

	if order.StudentID != UserData.UserID {
		lg.Info("[UserService] CreateReview now allowed. UserID: " + UserData.UserID + ". Order-UserID: " + order.StudentID)
		return "", custom_errors.ErrNotAllowed
	}

	if !response.IsFinal || order.Status != models.StatusSelected {
		lg.Info("[UserService] CreateReview bad status. Current order status: " + order.Status + " . Current response state: " + strconv.FormatBool(response.IsFinal))
		return "", custom_errors.ErrorBadStatus
	}

	if ReviewRequest.Rating < 0 || ReviewRequest.Rating > 5 {
		lg.Info("[UserService] Rating bad diapozon: " + strconv.Itoa(ReviewRequest.Rating))
		return "", custom_errors.ErrNotAllowed
	}

	tutor, err := u.GRPCClient.GetTutor(context.Background(), response.TutorID)

	if err != nil {
		lg.Error("[UserService] CreateReview error GetTutor: " + err.Error())
		return "", custom_errors.ErrorGetUser
	}

	currentTimestamp := time.Now()

	TimeDiff := currentTimestamp.Sub(response.CreatedAt)

	if TimeDiff < 1*time.Second {
		lg.Info("[UserService] Review time bad. Diff: " + TimeDiff.String())
		return "", custom_errors.ErrorLowTimeFromResponse
	}

	reviewID, err := u.GRPCClient.CreateReview(context.Background(), order.ID, response.TutorID, ReviewRequest.Comment, ReviewRequest.Rating)

	if err != nil {
		lg.Error("[UserService] CreateReview error CreateReview: " + err.Error())
		return "", custom_errors.ErrorServiceError
	}

	reviewToBroker := &models.ReviewToBroker{
		ReviewID:        reviewID,
		ResponseID:      response.ID,
		OrderID:         order.ID,
		OrderName:       order.Title,
		TutorTelegramID: tutor.TelegramID,
	}

	err = u.ProducerBroker.Publish("new_review", reviewToBroker)

	if err != nil {
		lg.Error("[UserService] CreateReview Error publishing order: " + err.Error())
	}

	err = u.orderRepository.SetOrderStatus(models.StatusClosed, order.ID)

	if err != nil {
		lg.Error("[UserService] CreateReview error SetOrderStatus: " + err.Error())
	}

	return reviewID, nil
}

func (u *UserService) GetReviewsByTutor(tutorID string, UserData models.UserData) ([]models.Review, error) {
	reviews, err := u.GRPCClient.GetReviewsByTutor(context.Background(), tutorID)
	if err != nil {
		lg.Error("[UserService] GetReviewsByTutor error GetReviewsByTutor: " + err.Error())
		return nil, err
	}

	return reviews, nil
}

func (u *UserService) GetReviewsByID(reviewID string, UserData models.UserData) (*models.Review, error) {
	review, err := u.GRPCClient.GetReviewsByID(context.Background(), reviewID)
	if err != nil {
		lg.Error("[UserService] GetReviewsByID error GetReviewsByID: " + err.Error())

		return nil, err
	}

	return review, nil
}

func (u *UserService) GetTutorInfoById(tutorID string, UserData models.UserData) (*models.TutorDetails, error) {
	isOwn := UserData.UserID == tutorID

	TutorDetails, err := u.GRPCClient.GetTutorInfoById(context.Background(), tutorID, isOwn)

	if err != nil {
		lg.Error("[UserService] GetTutorInfoById error GetTutorInfoById: " + err.Error())
		return nil, err
	}

	return TutorDetails, nil
}

func (u *UserService) ChangeTutorActive(isActive bool, UserData models.UserData) (bool, error) {
	lg.Info(fmt.Sprintf("[UserService] ChangeTutorActive. isActive: %v. UserID: %v. Role: %v", isActive, UserData.UserID, UserData.Role))

	isOk, err := u.GRPCClient.ChangeTutorActive(context.Background(), UserData.UserID, isActive)

	if err != nil {
		lg.Error(fmt.Sprintf("[UserService] ChangeTutorActive failed. isActive: %v. UserID: %v. Role: %v. Error: %v", isActive, UserData.UserID, UserData.Role, err.Error()))
		return false, err
	}

	lg.Info(fmt.Sprintf("[UserService] ChangeTutorActive success. isActive: %v. UserID: %v. Role: %v", isActive, UserData.UserID, UserData.Role))
	return isOk, nil
}

func (u *UserService) UpdateTutorName(tutorID, name string, UserData models.UserData) error {
	if utils.ContainsBadWords(name) {
		return custom_errors.ErrorBanWords
	}

	isOk, err := u.GRPCClient.UpdateNameTutor(context.Background(), tutorID, name)

	if err != nil || !isOk {
		return errors.New("cannot update tutor name")
	}

	return nil
}

func (u *UserService) SetReviewActive(reviewID string, UserData models.UserData) error {
	isOk, err := u.GRPCClient.SetActiveToReview(context.Background(), reviewID)

	if err != nil || !isOk {
		return errors.New("cannot update tutor name")
	}

	return nil
}

func (u *UserService) BanUser(banUser models.BanUser, UserData models.UserData) error {
	isOk, err := u.GRPCClient.BanUser(context.Background(), banUser.TelegramID, banUser.IsBan)

	if err != nil || !isOk {
		return err
	}

	return nil
}

func (u *UserService) GetUserById(userID string, UserData models.UserData) (*models.User, error) {
	return u.GRPCClient.GetUser(context.Background(), userID)
}
