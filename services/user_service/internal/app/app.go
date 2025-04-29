package app

import (
	"fmt"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/controllers"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type App struct {
	cfg        *config.Config
	repository repository.UserRepository

	userService     service.UserServiceInt
	userControllers *controllers.UserControllers
}

func NewApp(cfg *config.Config) *App {
	userRepo := repository.NewRepository(cfg.DBConfig)

	userServic := service.NewUserService(userRepo)
	userController := controllers.NewUserControllers(userServic)
	return &App{
		repository:      userRepo,
		userService:     userServic,
		userControllers: userController,
		cfg:             cfg,
	}
}

func (a *App) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.cfg.ServerPort))

	if err != nil {
		log.Fatal("server failed listen port " + a.cfg.ServerPort + " Error: " + err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, a.userControllers)

	log.Printf("server listening : %s", a.cfg.ServerPort)

	if err := s.Serve(lis); err != nil {
		log.Fatal("server failed!" + a.cfg.ServerPort + " Error: " + err.Error())
	}
}
