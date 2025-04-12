package repository

import (
	"github.com/randnull/Lessons/internal/models"
)

type OrderRepository interface {
	CreateOrder(order *models.NewOrder, studentID string, telegramID int64) (*models.OrderToBrokerWithID, error)
	GetOrderByID(orderID string) (*models.OrderDetails, error) // GetByID
	GetOrderByIdTutor(orderID string, tutorID string) (*models.OrderDetailsTutor, error)

	GetOrders() ([]*models.Order, error) // GetAllOrders
	GetOrdersPagination(limit int, offset int) ([]*models.Order, int, error)
	GetStudentOrders(studentID string) ([]*models.Order, error)

	GetStudentOrdersPagination(limit int, offset int, studentID string) ([]*models.Order, int, error)

	UpdateOrder(orderID string, order *models.UpdateOrder, studentID string) error
	GetUserByOrder(orderID string) (string, error)
	DeleteOrder(id string) error

	GetTutorsResponses(tutorID string) ([]models.Response, error)

	CheckOrderByStudentID(orderID string, studentID string) (bool, error)

	CheckResponseExist(TutorID, OrderID string) bool
	CreateResponse(orderID string, response *models.NewResponseModel, Tutor *models.Tutor, username string) (string, error)
	GetResponseById(ResponseID string, studentID string) (*models.ResponseDB, error)

	SetTutorToOrder(response *models.ResponseDB, UserData models.UserData) error
}
