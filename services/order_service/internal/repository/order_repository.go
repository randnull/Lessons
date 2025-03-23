package repository

import (
	"github.com/randnull/Lessons/internal/models"
)

type OrderRepository interface {
	CreateOrder(order *models.NewOrder, studentID string, telegramID int64) (*models.OrderToBrokerWithID, error)
	GetByID(id string, studentID string) (*models.OrderDetails, error)
	GetOrderByIdTutor(id string, tutorID string) (*models.OrderDetailsTutor, error)
	GetAllOrders(studentID string) ([]*models.Order, error)
	UpdateOrder(orderID string, order *models.UpdateOrder, studentID string) error
	GetUserByOrder(orderID string) (string, error)
	GetAllUsersOrders(studentID string) ([]*models.Order, error)
	DeleteOrder(id string, studentID string) error
	CreateResponse(orderID string, response *models.NewResponseModel, Tutor *models.User, username string) (string, error)
	GetResponseById(ResponseID string, studentID string) (*models.ResponseDB, error)
	CheckOrderByStudentID(orderID string, studentID string) (bool, error)
	SetTutorToOrder(responseID string, UserData models.UserData) error
}
