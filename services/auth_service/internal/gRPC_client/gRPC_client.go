package gRPC_client

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/pkg/custom_errors"
	pb "github.com/randnull/Lessons/pkg/gRPC"
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

	connectionLink := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)

	conn, err := grpc.Dial(connectionLink, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatal(err)
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

func (g *GRPCClient) GetUserByTelegramID(ctx context.Context, telegramID int64, role string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	userPB, err := g.client.GetUserByTelegramId(ctx, &pb.GetByTelegramId{Id: telegramID, Role: role})

	if err != nil {
		return nil, custom_errors.ErrorGetUser
	}

	return &models.User{
		Id:       userPB.Id,
		Name:     userPB.Name,
		IsBanned: userPB.IsBanned,
	}, nil
}

func (g *GRPCClient) CreateUser(ctx context.Context, user *models.NewUser) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	UserRPC, err := g.client.CreateUser(ctx, &pb.CreateUserRequest{Name: user.Name, TelegramId: user.TelegramID, Role: user.Role})

	if err != nil {
		return "", custom_errors.ErrorCreateUser
	}

	return UserRPC.Id, nil
}
