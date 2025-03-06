package gRPC_client

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/custom_errors"
	pb "github.com/randnull/Lessons/internal/gRPC"
	"github.com/randnull/Lessons/internal/models"
	"google.golang.org/grpc"
	"log"
	"time"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.PostsServiceClient
}

func NewGRPCClient() *GRPCClient {
	fmt.Println("Waiting connection")

	// FATAL!!!! ОЖИДАЕТ Connection до КОНЦА!! СРОЧНО ИСПРАВИТЬ
	conn, err := grpc.Dial("lessons-user-service:2000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("Can't establish connect with gRPC. Fatal Error")
	}
	client := pb.NewPostsServiceClient(conn)

	log.Println("Connection with gRPC: ok")
	return &GRPCClient{
		conn:   conn,
		client: client,
	}
}

func (g GRPCClient) Close() {
	err := g.conn.Close()
	if err != nil {
		log.Printf("error with close connection")
	}
}

func (g GRPCClient) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	userPB, err := g.client.GetUserById(ctx, &pb.GetById{Id: userID})

	if err != nil {
		return nil, custom_errors.ErrorGetUser
	}

	return &models.User{
		Id:   userPB.Id,
		Name: userPB.Name,
	}, nil
}

func (g GRPCClient) CreateUser(ctx context.Context, user *models.CreateUser) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	userID, err := g.client.CreateUser(ctx, &pb.CreateUserRequest{Name: user.Name, TelegramId: user.TelegramId})

	if err != nil {
		return "", custom_errors.ErrorCreateUser
	}

	return userID.Id, nil
}
