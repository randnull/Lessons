package app

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/controllers"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/scheduler"
	"github.com/randnull/Lessons/internal/service"
	pb "github.com/randnull/Lessons/pkg/gRPC"
	lg "github.com/randnull/Lessons/pkg/logger"
	"github.com/randnull/Lessons/pkg/rabbitmq"
	"google.golang.org/grpc"

	"net"
)

type App struct {
	cfg        *config.Config
	repository repository.UserRepository

	userService     service.UserServiceInt
	userControllers *controllers.UserControllers

	scheduler scheduler.Scheduler
}

func NewApp(cfg *config.Config) *App {
	userRepo := repository.NewRepository(cfg.DBConfig)

	BrokerProducer := rabbitmq.NewRabbitMQ(cfg.MQConfig)

	userServic := service.NewUserService(userRepo)
	userController := controllers.NewUserControllers(userServic)

	regularScheduler := *scheduler.NewScheduler(&cfg.SchedulerConfig, userRepo, BrokerProducer)

	return &App{
		repository:      userRepo,
		userService:     userServic,
		userControllers: userController,
		scheduler:       regularScheduler,
		cfg:             cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.scheduler.RunResponseChecker(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.cfg.ServerPort))

	if err != nil {
		lg.Info("server failed listen port " + a.cfg.ServerPort + " Error: " + err.Error())
		return err
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, a.userControllers)

	go func() {
		err = s.Serve(lis)
		if err != nil {
			lg.Info("server failed serve " + a.cfg.ServerPort + " Error: " + err.Error())
		}
	}()

	lg.Info("gRPC server listening: " + a.cfg.ServerPort)

	<-ctx.Done()

	s.GracefulStop()
	lg.Info("gRPC successfully shutdown")

	return nil
}
