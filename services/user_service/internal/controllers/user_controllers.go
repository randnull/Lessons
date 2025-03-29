package controllers

import (
	"context"
	"fmt"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
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

func (s *UserControllers) GetAllTutorsPagination(ctx context.Context, in *pb.GetAllTutorsPaginationRequest) (*pb.GetAllResponse, error) {
	//page, err := strconv.Atoi(in.Page)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//size, err := strconv.Atoi(in.Size)
	//
	//if err != nil {
	//	return nil, err
	//}

	users, err := s.UserService.GetTutorsPagination(int(in.Page), int(in.Size))

	if err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{
		Users: users,
	}, nil
}

func (s *UserControllers) UpdateBioTutor(ctx context.Context, in *pb.UpdateBioRequest) (*pb.UpdateBioResponse, error) {
	err := s.UserService.UpdateBioTutor(in.Id, in.Bio)

	if err != nil {
		return &pb.UpdateBioResponse{Success: false}, err
	}

	return &pb.UpdateBioResponse{Success: true}, nil
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
