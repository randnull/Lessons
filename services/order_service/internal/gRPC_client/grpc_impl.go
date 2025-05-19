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
	connectionLink := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)

	conn, err := grpc.Dial(connectionLink, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("error connection to grpc: " + err.Error())
	}

	client := pb.NewUserServiceClient(conn)
	return &GRPCClient{
		conn:   conn,
		client: client,
	}
}

func (g *GRPCClient) GetUser(ctx context.Context, userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userPB, err := g.client.GetUserById(ctx, &pb.GetById{Id: userID})
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:         userPB.Id,
		TelegramID: userPB.TelegramId,
		Name:       userPB.Name,
		Role:       userPB.Role,
		IsBanned:   userPB.IsBanned,
	}, nil
}

func (g *GRPCClient) GetStudent(ctx context.Context, userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userPB, err := g.client.GetStudentById(ctx, &pb.GetById{Id: userID})
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:         userPB.Id,
		TelegramID: userPB.TelegramId,
		Name:       userPB.Name,
	}, nil
}

func (g *GRPCClient) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userPB, err := g.client.GetUserByTelegramId(ctx, &pb.GetByTelegramId{Id: telegramID})
	if err != nil {
		return nil, err
	}

	return &models.User{
		Id:         userPB.Id,
		TelegramID: userPB.TelegramId,
		Name:       userPB.Name,
	}, nil
}

func (g *GRPCClient) GetAllUsers(ctx context.Context) (*pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	usersPB, err := g.client.GetAllUsers(ctx, &pb.GetAllRequest{})

	if err != nil {
		return nil, err
	}

	return usersPB, nil
}

func (g *GRPCClient) UpdateBioTutor(ctx context.Context, bio string, tutorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := g.client.UpdateBioTutor(ctx, &pb.UpdateBioRequest{
		Id:  tutorID,
		Bio: bio,
	})
	if err != nil {
		return false, err
	}

	return status.Success, nil
}

func (g *GRPCClient) GetTutor(ctx context.Context, tutorID string) (*models.Tutor, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tutorPB, err := g.client.GetTutorById(ctx, &pb.GetById{Id: tutorID})
	if err != nil {
		return nil, err
	}

	return &models.Tutor{
		Id:         tutorPB.User.Id,
		TelegramID: tutorPB.User.TelegramId,
		Name:       tutorPB.User.Name,
		Tags:       tutorPB.Tags,
	}, nil
}

func (g *GRPCClient) GetTutorsPagination(ctx context.Context, page int, size int, tag string) (*pb.GetTutorsPaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	usersPB, err := g.client.GetAllTutorsPagination(ctx, &pb.GetAllTutorsPaginationRequest{
		Page: int32(page),
		Size: int32(size),
		Tag:  tag,
	})
	if err != nil {
		return nil, err
	}

	return usersPB, nil
}

func (g *GRPCClient) UpdateTagsTutor(ctx context.Context, tags []string, tutorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.UpdateTags(ctx, &pb.UpdateTagsRequest{
		TutorId: tutorID,
		Tags:    tags,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (g *GRPCClient) CreateReview(ctx context.Context, orderID string, tutorID string, comment string, rating int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.CreateReview(ctx, &pb.CreateReviewRequest{
		TutorId: tutorID,
		OrderId: orderID,
		Rating:  int32(rating),
		Comment: comment,
	})
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}

func (g *GRPCClient) GetReviewsByTutor(ctx context.Context, tutorID string) ([]models.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.GetReviews(ctx, &pb.GetReviewsRequest{TutorId: tutorID})

	if err != nil {
		return nil, err
	}

	var reviews []models.Review
	for _, r := range resp.Reviews {
		reviews = append(reviews, models.Review{
			ID:        r.Id,
			TutorID:   r.TutorId,
			OrderID:   r.OrderId,
			Rating:    int(r.Rating),
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt.AsTime(),
		})
	}

	return reviews, nil
}

func (g *GRPCClient) GetReviewsByID(ctx context.Context, reviewID string) (*models.Review, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.GetReview(ctx, &pb.GetReviewRequest{Id: reviewID})

	if err != nil {
		return nil, err
	}

	if resp.IsActive {
		return nil, custom_errors.ErrorNotActiveReview
	}

	return &models.Review{
		ID:        resp.Id,
		TutorID:   resp.TutorId,
		OrderID:   resp.OrderId,
		Rating:    int(resp.Rating),
		Comment:   resp.Comment,
		CreatedAt: resp.CreatedAt.AsTime(),
	}, nil
}

func (g *GRPCClient) GetTutorInfoById(ctx context.Context, tutorID string, isOwn bool) (*models.TutorDetails, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.GetTutorInfoById(ctx, &pb.GetById{Id: tutorID})
	if err != nil {
		return nil, err
	}

	var reviews []models.Review
	for _, r := range resp.Review {
		if isOwn || r.IsActive {
			reviews = append(reviews, models.Review{
				ID:        r.Id,
				TutorID:   r.TutorId,
				OrderID:   r.OrderId,
				Rating:    int(r.Rating),
				Comment:   r.Comment,
				IsActive:  r.IsActive,
				CreatedAt: r.CreatedAt.AsTime(),
			})
		}
	}

	return &models.TutorDetails{
		Tutor: models.User{
			Id:         resp.User.Id,
			Name:       resp.User.Name,
			TelegramID: resp.User.TelegramId,
			Role:       resp.User.Role,
			IsBanned:   resp.User.IsBanned,
		},
		Bio:           resp.Bio,
		ResponseCount: resp.ResponseCount,
		Reviews:       reviews,
		Tags:          resp.Tags,
		Rating:        resp.Rating,
		IsActive:      resp.IsActive,
		CreatedAt:     resp.CreatedAt.AsTime(),
	}, nil
}

func (g *GRPCClient) ChangeTutorActive(ctx context.Context, tutorID string, active bool) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.ChangeTutorActive(ctx, &pb.SetActiveTutorById{
		Id:     tutorID,
		Active: active,
	})

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (g *GRPCClient) CreateNewResponse(ctx context.Context, tutorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.CreateNewResponse(ctx, &pb.CreateResponseRequest{
		TutorId: tutorID,
	})
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (g *GRPCClient) UpdateNameTutor(ctx context.Context, tutorID, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status, err := g.client.ChangeTutorName(ctx, &pb.ChangeNameRequest{
		Id:   tutorID,
		Name: name,
	})
	if err != nil {
		return false, err
	}

	return status.Success, nil
}

func (g *GRPCClient) SetActiveToReview(ctx context.Context, reviewID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.SetReviewActive(ctx, &pb.SetReviewsActiveRequest{
		ReviewId: reviewID,
	})

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (g *GRPCClient) BanUser(ctx context.Context, telegramID int64, isBanned bool) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := g.client.BanUser(ctx, &pb.BanUserRequest{
		TelegramId: telegramID,
		IsBanned:   isBanned,
	})

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (g *GRPCClient) Close() error {
	return g.conn.Close()
}
