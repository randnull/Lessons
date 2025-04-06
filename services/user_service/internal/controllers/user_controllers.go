package controllers

import (
	"context"
	"fmt"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
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
	fmt.Println(in)
	userID, err := s.UserService.CreateUser(models.CreateUser{Name: in.Name, TelegramId: in.TelegramId, Role: in.Role})
	if err != nil {
		return nil, err
	}

	userPB := &pb.CreateResponse{
		Id: userID,
	}

	return userPB, nil
}

//func (s *UserControllers) GetUserByTelegramId(ctx context.Context, in *pb.GetByTelegramId) (*pb.User, error) {
//	user, err := s.UserService.GetUserByTelegramId(in.Id)
//
//	if err != nil {
//		return nil, err
//	}
//
//	userPB := &pb.User{
//		Id:   user.Id,
//		Name: user.Name,
//	}
//	return userPB, nil
//}

func (s *UserControllers) GetUserById(ctx context.Context, in *pb.GetById) (*pb.User, error) {
	user, err := s.UserService.GetUserById(in.Id)

	if err != nil {
		return nil, err
	}

	userPB := &pb.User{
		Id:         user.Id,
		TelegramId: user.TelegramID,
		Name:       user.Name,
	}
	return userPB, nil
}

func (s *UserControllers) GetStudentById(ctx context.Context, in *pb.GetById) (*pb.User, error) {
	user, err := s.UserService.GetStudentById(in.Id)

	if err != nil {
		return nil, err
	}

	userPB := &pb.User{
		Id:         user.Id,
		TelegramId: user.TelegramID,
		Name:       user.Name,
	}
	return userPB, nil
}

func (s *UserControllers) GetAllUsers(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := s.UserService.GetTutors()

	if err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{
		Users: users,
	}, nil
}

func (s *UserControllers) GetAllTutorsPagination(ctx context.Context, in *pb.GetAllTutorsPaginationRequest) (*pb.GetTutorsPaginationResponse, error) {
	users, err := s.UserService.GetTutorsPagination(int(in.Page), int(in.Size))

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserControllers) UpdateBioTutor(ctx context.Context, in *pb.UpdateBioRequest) (*pb.UpdateBioResponse, error) {
	err := s.UserService.UpdateBioTutor(in.Id, in.Bio)

	if err != nil {
		return &pb.UpdateBioResponse{Success: false}, err
	}

	return &pb.UpdateBioResponse{Success: true}, nil
}

func (s *UserControllers) UpdateTags(ctx context.Context, in *pb.UpdateTagsRequest) (*pb.Success, error) {
	err := s.UserService.UpdateTutorTags(in.TutorId, in.Tags)
	if err != nil {
		return &pb.Success{Success: false}, err
	}

	return &pb.Success{Success: true}, nil
}

func (s *UserControllers) GetTutorById(ctx context.Context, in *pb.GetById) (*pb.Tutor, error) {
	tutor, err := s.UserService.GetTutorById(in.Id)

	if err != nil {
		return nil, err
	}

	userPB := &pb.User{
		Id:         tutor.Id,
		TelegramId: tutor.TelegramID,
		Name:       tutor.Name,
	}

	return &pb.Tutor{
		Bio:  tutor.Bio,
		User: userPB,
	}, nil
}

func (s *UserControllers) ChangeTutorActive(ctx context.Context, in *pb.SetActiveTutorById) (*pb.Success, error) {
	err := s.UserService.ChangeTutorActive(in.Id, in.Active)
	if err != nil {
		return &pb.Success{Success: false}, err
	}
	return &pb.Success{Success: true}, nil
}

func (s *UserControllers) GetTutorInfoById(ctx context.Context, in *pb.GetById) (*pb.TutorDetails, error) {
	tutor, err := s.UserService.GetTutorInfoById(in.Id)
	if err != nil {
		return nil, err
	}

	var reviews []*pb.Review
	for _, r := range tutor.Reviews {
		reviews = append(reviews, &pb.Review{
			Id:        r.ID,
			TutorId:   r.TutorID,
			StudentId: r.StudentID,
			Rating:    int32(r.Rating),
			Comment:   r.Comment,
			CreatedAt: timestamppb.New(r.CreatedAt),
		})
	}

	return &pb.TutorDetails{
		Tutor: &pb.Tutor{
			User: &pb.User{
				Id:         tutor.Tutor.Id,
				Name:       tutor.Tutor.Name,
				TelegramId: tutor.Tutor.TelegramID,
			},
			Bio: tutor.Tutor.Bio,
		},
		Bio:    tutor.Tutor.Bio,
		Tags:   tutor.Tags,
		Review: reviews,
	}, nil
}

func (s *UserControllers) CreateReview(ctx context.Context, in *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	reviewID, err := s.UserService.CreateReview(in.TutorId, in.StudentId, int(in.Rating), in.Comment)
	if err != nil {
		return nil, err
	}

	return &pb.CreateReviewResponse{
		Id: reviewID,
	}, nil
}

func (s *UserControllers) GetReview(ctx context.Context, in *pb.GetReviewRequest) (*pb.Review, error) {
	review, err := s.UserService.GetReviewById(in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Review{
		Id:        review.ID,
		TutorId:   review.TutorID,
		StudentId: review.StudentID,
		Rating:    int32(review.Rating),
		Comment:   review.Comment,
		CreatedAt: timestamppb.New(review.CreatedAt),
	}, nil
}

func (s *UserControllers) GetReviews(ctx context.Context, in *pb.GetReviewsRequest) (*pb.GetReviewsResponse, error) {
	reviews, err := s.UserService.GetReviews(in.TutorId)
	if err != nil {
		return nil, err
	}

	var pbReviews []*pb.Review
	for _, r := range reviews {
		pbReviews = append(pbReviews, &pb.Review{
			Id:        r.ID,
			TutorId:   r.TutorID,
			StudentId: r.StudentID,
			Rating:    int32(r.Rating),
			Comment:   r.Comment,
			CreatedAt: timestamppb.New(r.CreatedAt),
		})
	}

	return &pb.GetReviewsResponse{
		Reviews: pbReviews,
	}, nil
}
