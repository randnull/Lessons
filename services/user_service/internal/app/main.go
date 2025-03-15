package app

import (
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
	userRepo := repository.NewRepository(cfg.DBConfig) //cfg.DBConfig
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
	lis, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal("server failed")
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, a.userControllers)

	log.Printf("server listening : %s", "2000")

	if err := s.Serve(lis); err != nil {
		log.Fatal("server failed")
	}
}
