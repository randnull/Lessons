package redis

import (
	"context"
	"github.com/randnull/Lessons/internal/models"
)

type RedisInterface interface {
	SaveNewOrder(order *models.Order) error
	GetOrderHSET(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrderHSET(orderID string, updateOrder *models.UpdateOrder) error
	SaveNewResponse(response *models.ResponseDB) error
	GetAllResponses(ctx context.Context, orderID string) ([]models.Response, error)
	GetResponseById(ctx context.Context, responseID string) (*models.ResponseDB, error)
	GetShortResponseById(ctx context.Context, responseID string) (*models.Response, error)
	DeleteOrderHSET(orderID string) error
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
	GetAllUsersOrders(ctx context.Context, userID string) ([]*models.Order, error)
	AddOrder(order *models.OrderDetails) error
	AddResponses(response []*models.Response) error
	AddResponse(response *models.Response) error
}
