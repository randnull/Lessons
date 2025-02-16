package gRPC_client

import (
	"context"
	"fmt"
	pb "github.com/randnull/Lessons/internal/gRPC"
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

func (g GRPCClient) GetUser(ctx context.Context, userID string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return g.client.GetUserById(ctx, &pb.GetById{Id: userID})
}
