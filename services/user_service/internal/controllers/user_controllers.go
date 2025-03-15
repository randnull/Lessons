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
		Id:   user.Id,
		Name: user.Name,
	}
	return userPB, nil
}

func (s *UserControllers) GetAllUsers(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := s.UserService.GetAllUsers()

	if err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{
		Users: users,
	}, nil
}
