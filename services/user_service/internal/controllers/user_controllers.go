package controllers

import (
	"context"
	"fmt"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/service"
	"log"
)

type UserControllers struct {
	UserService service.UserServiceInt
	pb.UnimplementedPostsServiceServer
}

func NewUserControllers(userService service.UserServiceInt) *UserControllers {
	return &UserControllers{
		UserService: userService,
	}
}

func (s *UserControllers) GetUserById(ctx context.Context, in *pb.GetById) (*pb.User, error) {
	//post, err := s.repo.GetById(in.Id, int(in.UserId))
	user, err := s.UserService.GetUserById(in.Id)
	fmt.Println(in.Id)
	if err != nil {
		log.Fatal("aaaaa")
	}
	//if err != nil {
	//	status = 1
	//	postAnswer := &pb.Post{
	//		Id:        "",
	//		UserId:    0,
	//		Title:     "",
	//		Body:      "",
	//		CreatedAt: timestamppb.New(time.Now()),
	//		Status:    int64(status),
	//	}
	//	return postAnswer, nil
	//}

	userPB := &pb.User{
		Id:   user.UserId,
		Name: user.Name,
	}

	return userPB, nil
}
