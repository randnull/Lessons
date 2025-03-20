package gRPC_client

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/custom_errors"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"google.golang.org/grpc"
	"log"
	"time"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewGRPCClient(cfg config.GRPCConfig) *GRPCClient {
	fmt.Println("Waiting connection")
	//
	//var retryPolicy = `{
	//        "methodConfig": [{
	//            // config per method or all methods under service
	//            "name": [{"service": "grpc.examples.echo.Echo"}],
	//
	//            "retryPolicy": {
	//                "MaxAttempts": 4,
	//                "InitialBackoff": ".01s",
	//                "MaxBackoff": ".01s",
	//                "BackoffMultiplier": 1.0,
	//                // this value is grpc code
	//                "RetryableStatusCodes": [ "UNAVAILABLE" ]
	//            }
	//        }]
	//    }`

	connection_link := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Println(connection_link)
	// FATAL!!!! ОЖИДАЕТ Connection до КОНЦА!! СРОЧНО ИСПРАВИТЬ
	conn, err := grpc.Dial(connection_link, grpc.WithInsecure(), grpc.WithBlock())
	//conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Can't establish connect with gRPC. Fatal Error")
	}
	client := pb.NewUserServiceClient(conn)

	log.Println("Connection with gRPC: ok")
	return &GRPCClient{
		conn:   conn,
		client: client,
	}
}

func (g *GRPCClient) Close() {
	err := g.conn.Close()
	if err != nil {
		log.Printf("error with close connection")
	}
}

func (g *GRPCClient) GetUser(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userPB, err := g.client.GetUserById(ctx, &pb.GetById{Id: userID})

	if err != nil {
		return nil, custom_errors.ErrorGetUser
	}

	return &models.User{
		Id:         userPB.Id,
		TelegramID: userPB.TelegramId,
		Name:       userPB.Name,
	}, nil
}

func (g *GRPCClient) GetStudent(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userPB, err := g.client.GetStudentById(ctx, &pb.GetById{Id: userID})

	if err != nil {
		return nil, custom_errors.ErrorGetUser
	}

	return &models.User{
		Id:         userPB.Id,
		TelegramID: userPB.TelegramId,
		Name:       userPB.Name,
	}, nil
}

func (g *GRPCClient) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userPB, err := g.client.GetUserByTelegramId(ctx, &pb.GetByTelegramId{Id: telegramID})

	if err != nil {
		return nil, custom_errors.ErrorGetUser
	}

	return &models.User{
		Id:         userPB.Id,
		TelegramID: userPB.TelegramId,
		Name:       userPB.Name,
	}, nil
}

func (g *GRPCClient) GetAllUsers() (*pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersPB, err := g.client.GetAllUsers(ctx, &pb.GetAllRequest{})

	if err != nil {
		return nil, custom_errors.ErrorGetUser
	}

	return usersPB, nil
}
