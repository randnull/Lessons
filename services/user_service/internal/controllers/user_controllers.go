package controllers

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	pb "github.com/randnull/Lessons/pkg/gRPC"
	custom_logger "github.com/randnull/Lessons/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserControllers struct {
	UserService service.UserServiceInt
	pb.UnimplementedUserServiceServer
}

func NewUserControllers(userService service.UserServiceInt) *UserControllers {
	return &UserControllers{
		UserService: userService,
	}
}

func (s *UserControllers) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateResponse, error) {
	userID, err := s.UserService.CreateUser(models.CreateUser{Name: in.Name, TelegramId: in.TelegramId, Role: in.Role})
	if err != nil {
		custom_logger.Error("CreateUser error. UserID: " + in.Name + ", Role: " + in.Role + ", TelegramID: " + fmt.Sprint(in.TelegramId) + "Error: " + err.Error())
		return nil, err
	}

	userPB := &pb.CreateResponse{
		Id: userID,
	}

	custom_logger.Info("CreateUser Created. UserID: " + in.Name + ", Role: " + in.Role + ", TelegramID: " + fmt.Sprint(in.TelegramId) + "ID: " + userID)

	return userPB, nil
}

func (s *UserControllers) GetUserByTelegramId(ctx context.Context, in *pb.GetByTelegramId) (*pb.User, error) {
	custom_logger.Info("GetUserByTelegramId called. UserTelegramID: " + fmt.Sprint(in.Id))

	user, err := s.UserService.GetUserByTelegramId(in.Id, in.Role)

	if err != nil {
		custom_logger.Error("GetUserByTelegramId failed. UserTelegramID: " + fmt.Sprint(in.Id) + "Error: " + err.Error())
		return nil, err
	}

	userPB := &pb.User{
		Id:   user.Id,
		Name: user.Name,
	}

	custom_logger.Info("GetUserByTelegramId success. UserTelegramID: " + fmt.Sprint(in.Id) + " UserID: " + userPB.Id)

	return userPB, nil
}

func (s *UserControllers) GetStudentById(ctx context.Context, in *pb.GetById) (*pb.User, error) {
	custom_logger.Info("GetStudentById called. UserID: " + in.Id)

	user, err := s.UserService.GetStudentById(in.Id)

	if err != nil {
		custom_logger.Error("GetStudentById failed. UserID: " + in.Id + "Error: " + err.Error())
		return nil, err
	}

	userPB := &pb.User{
		Id:         user.Id,
		TelegramId: user.TelegramID,
		Name:       user.Name,
	}

	custom_logger.Info("GetStudentById success. UserID: " + in.Id)

	return userPB, nil
}

func (s *UserControllers) GetTutorById(ctx context.Context, in *pb.GetById) (*pb.Tutor, error) {
	custom_logger.Info("GetTutorById called. UserID: " + in.Id)

	tutor, err := s.UserService.GetTutorById(in.Id)

	if err != nil {
		custom_logger.Error("GetTutorById failed. UserID: " + in.Id + "Error: " + err.Error())
		return nil, err
	}

	userPB := &pb.User{
		Id:         tutor.Id,
		TelegramId: tutor.TelegramID,
		Name:       tutor.Name,
	}

	custom_logger.Info("GetTutorById success. UserID: " + in.Id)

	return &pb.Tutor{
		User: userPB,
		Tags: tutor.Tags,
	}, nil
}

func (s *UserControllers) GetAllUsers(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	custom_logger.Info("GetAllUsers called")

	users, err := s.UserService.GetAllUsers()

	if err != nil {
		custom_logger.Error("GetAllUsers failed. Error: " + err.Error())
		return nil, err
	}

	custom_logger.Info("GetAllUsers success")

	return &pb.GetAllResponse{
		Users: users,
	}, nil
}

func (s *UserControllers) GetAllTutorsPagination(ctx context.Context, in *pb.GetAllTutorsPaginationRequest) (*pb.GetTutorsPaginationResponse, error) {
	custom_logger.Info("GetAllTutorsPagination called.")

	tutors, err := s.UserService.GetTutorsPagination(int(in.Page), int(in.Size), in.Tag)

	if err != nil {
		custom_logger.Error("GetAllTutorsPagination failed. Error: " + err.Error())
		return nil, err
	}

	custom_logger.Info("GetAllTutorsPagination success")

	return tutors, nil
}

func (s *UserControllers) UpdateBioTutor(ctx context.Context, in *pb.UpdateBioRequest) (*pb.Success, error) {
	custom_logger.Info("UpdateBioTutor called. TutorID: " + in.Id)

	err := s.UserService.UpdateBioTutor(in.Id, in.Bio)

	if err != nil {
		custom_logger.Error("UpdateBioTutor failed. TutorID: " + in.Id + "Error" + err.Error())
		return &pb.Success{Success: false}, err
	}

	custom_logger.Info("UpdateBioTutor success")
	return &pb.Success{Success: true}, nil
}

func (s *UserControllers) UpdateTags(ctx context.Context, in *pb.UpdateTagsRequest) (*pb.Success, error) {
	custom_logger.Info("UpdateTags called. TutorID: " + in.TutorId)

	err := s.UserService.UpdateTutorTags(in.TutorId, in.Tags)
	if err != nil {
		custom_logger.Error("UpdateTags failed. TutorID: " + in.TutorId + "Error" + err.Error())
		return &pb.Success{Success: false}, err
	}

	custom_logger.Info("UpdateTags success. TutorID: " + in.TutorId)
	return &pb.Success{Success: true}, nil
}

func (s *UserControllers) ChangeTutorActive(ctx context.Context, in *pb.SetActiveTutorById) (*pb.Success, error) {
	custom_logger.Info("ChangeTutorActive called. TutorID: " + in.Id)

	err := s.UserService.ChangeTutorActive(in.Id, in.Active)

	if err != nil {
		custom_logger.Error("ChangeTutorActive failed. TutorID: " + in.Id + "Error" + err.Error())
		return &pb.Success{Success: false}, err
	}

	custom_logger.Info("ChangeTutorActive success. TutorID: " + in.Id)
	return &pb.Success{Success: true}, nil
}

func (s *UserControllers) GetTutorInfoById(ctx context.Context, in *pb.GetById) (*pb.TutorDetails, error) {
	custom_logger.Info("GetTutorInfoById called. TutorID: " + in.Id)

	tutor, err := s.UserService.GetTutorInfoById(in.Id)

	if err != nil {
		custom_logger.Error("GetTutorInfoById failed GetTutorInfoById. TutorID: " + in.Id + "Error" + err.Error())
		return nil, err
	}

	var reviews []*pb.Review
	for _, r := range tutor.Reviews {
		reviews = append(reviews, &pb.Review{
			Id:        r.ID,
			TutorId:   r.TutorID,
			OrderId:   r.OrderID,
			Rating:    int32(r.Rating),
			Comment:   r.Comment,
			IsActive:  r.IsActive,
			CreatedAt: timestamppb.New(r.CreatedAt),
		})
	}

	custom_logger.Info("GetTutorInfoById success. TutorID: " + in.Id)

	return &pb.TutorDetails{
		User: &pb.User{
			Id:         tutor.Tutor.Id,
			Name:       tutor.Tutor.Name,
			TelegramId: tutor.Tutor.TelegramID,
		},
		IsActive:      tutor.Tutor.IsActive,
		ResponseCount: tutor.Tutor.ResponseCount,
		Bio:           tutor.Tutor.Bio,
		Rating:        tutor.Tutor.Rating,
		Tags:          tutor.Tutor.Tags,
		Review:        reviews,
		CreatedAt:     timestamppb.New(tutor.Tutor.CreatedAt),
	}, nil
}

func (s *UserControllers) CreateReview(ctx context.Context, in *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	custom_logger.Info("CreateReview called. TutorID: " + in.TutorId + " OrderID: " + in.OrderId)

	reviewID, err := s.UserService.CreateReview(in.TutorId, in.OrderId, int(in.Rating), in.Comment)

	if err != nil {
		custom_logger.Error("CreateReview failed. TutorID: " + in.TutorId + " OrderID: " + in.OrderId)
		return nil, err
	}

	custom_logger.Info("CreateReview success. TutorID: " + in.TutorId + " OrderID: " + in.OrderId)

	return &pb.CreateReviewResponse{
		Id: reviewID,
	}, nil
}

func (s *UserControllers) GetReview(ctx context.Context, in *pb.GetReviewRequest) (*pb.Review, error) {
	custom_logger.Info("GetReview called. ReviewID: " + in.Id)

	review, err := s.UserService.GetReviewById(in.Id)

	if err != nil {
		custom_logger.Error("GetReview failed. ReviewID: " + in.Id)
		return nil, err
	}

	custom_logger.Info("GetReview success. ReviewID: " + in.Id)

	return &pb.Review{
		Id:        review.ID,
		TutorId:   review.TutorID,
		OrderId:   review.OrderID,
		Rating:    int32(review.Rating),
		Comment:   review.Comment,
		IsActive:  review.IsActive,
		CreatedAt: timestamppb.New(review.CreatedAt),
	}, nil
}

func (s *UserControllers) GetReviews(ctx context.Context, in *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	custom_logger.Info("GetReviews called. TutorID: " + in.TutorId)

	reviews, err := s.UserService.GetReviews(in.TutorId)

	if err != nil {
		custom_logger.Error("GetReviews failed. TutorID: " + in.TutorId)
		return nil, err
	}

	custom_logger.Info("GetReviews success. TutorID: " + in.TutorId)

	var pbReviews []*pb.Review
	for _, r := range reviews {
		pbReviews = append(pbReviews, &pb.Review{
			Id:        r.ID,
			TutorId:   r.TutorID,
			OrderId:   r.OrderID,
			Rating:    int32(r.Rating),
			Comment:   r.Comment,
			CreatedAt: timestamppb.New(r.CreatedAt),
		})
	}

	return &pb.GetReviewsResponse{
		Reviews: pbReviews,
	}, nil
}

func (s *UserControllers) CreateNewResponse(ctx context.Context, in *pb.CreateResponseRequest) (*pb.Success, error) {
	custom_logger.Info("CreateNewResponse called. TutorID: " + in.TutorId)

	err := s.UserService.CreateNewResponse(in.TutorId)

	if err != nil {
		custom_logger.Error("CreateNewResponse failed. TutorID: " + in.TutorId)
		return &pb.Success{Success: false}, err
	}

	custom_logger.Info("CreateNewResponse success. TutorID: " + in.TutorId)

	return &pb.Success{Success: true}, nil
}

func (s *UserControllers) AddResponsesToTutor(ctx context.Context, in *pb.AddResponseToTutorRequest) (*pb.AddResponseToTutorResponse, error) {
	custom_logger.Info("AddResponsesToTutor called. TutorID: " + fmt.Sprint(in.TutorId))

	responses, err := s.UserService.AddResponses(in.TutorId, int(in.ResponseCount))

	if err != nil {
		custom_logger.Error("AddResponsesToTutor failed. TutorID: " + fmt.Sprint(in.TutorId))
		return &pb.AddResponseToTutorResponse{
			ResponseCount: 0,
			Success:       false,
		}, err
	}

	custom_logger.Info("AddResponsesToTutor success. TutorID: " + fmt.Sprint(in.TutorId))

	return &pb.AddResponseToTutorResponse{
		ResponseCount: int32(responses),
		Success:       true,
	}, nil
}

func (s *UserControllers) ChangeTutorName(ctx context.Context, in *pb.ChangeNameRequest) (*pb.Success, error) {
	custom_logger.Info("ChangeTutorName called. TutorID: " + in.Id)

	err := s.UserService.UpdateNameTutor(in.Id, in.Name)

	if err != nil {
		custom_logger.Error("ChangeTutorName failed. TutorID: " + in.Id)
		return &pb.Success{
			Success: false,
		}, err
	}

	custom_logger.Info("ChangeTutorName success. TutorID: " + in.Id)

	return &pb.Success{
		Success: true,
	}, nil
}

func (s *UserControllers) SetReviewActive(ctx context.Context, in *pb.SetReviewsActiveRequest) (*pb.Success, error) {
	custom_logger.Info("SetReviewActive called. ReviewID: " + in.ReviewId)

	err := s.UserService.SetReviewActive(in.ReviewId)

	if err != nil {
		custom_logger.Error("SetReviewActive failed. ReviewID: " + in.ReviewId)
		return &pb.Success{
			Success: false,
		}, err
	}

	custom_logger.Info("SetReviewActive success. ReviewID: " + in.ReviewId)

	return &pb.Success{
		Success: true,
	}, nil
}

func (s *UserControllers) BanUser(ctx context.Context, in *pb.BanUserRequest) (*pb.Success, error) {
	custom_logger.Info(fmt.Sprintf("BanUser called. TelegramID: %v", in.TelegramId))

	err := s.UserService.BanUser(in.TelegramId, in.IsBanned)

	if err != nil {
		custom_logger.Error(fmt.Sprintf("BanUser failed. TelegramID: %v", in.TelegramId))
		return &pb.Success{
			Success: false,
		}, err
	}

	custom_logger.Info(fmt.Sprintf("BanUser success. TelegramID: %v", in.TelegramId))

	return &pb.Success{
		Success: true,
	}, nil
}

func (s *UserControllers) GetUserById(ctx context.Context, in *pb.GetById) (*pb.User, error) {
	custom_logger.Info("GetStudentById called. UserID: " + in.Id)

	user, err := s.UserService.GetUserById(in.Id)

	if err != nil {
		custom_logger.Error("GetStudentById failed. UserID: " + in.Id + "Error: " + err.Error())
		return nil, err
	}

	userPB := &pb.User{
		Id:         user.Id,
		TelegramId: user.TelegramID,
		Name:       user.Name,
		Role:       user.Role,
		IsBanned:   user.IsBanned,
	}

	custom_logger.Info("GetStudentById success. UserID: " + in.Id)

	return userPB, nil
}
