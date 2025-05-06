package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/randnull/Lessons/internal/models"
)

type OrderRepository interface {
	// Order Edit
	CreateOrder(NewOrder *models.CreateOrder) (string, error)
	UpdateOrder(orderID string, order *models.UpdateOrder) error
	DeleteOrder(id string) error
	SetOrderStatus(status string, orderID string) error
	SetTutorToOrder(response *models.ResponseDB, UserData models.UserData) error

	// Order Getting
	GetOrderByID(orderID string) (*models.Order, error)
	GetOrders() ([]*models.Order, error) // GetAllOrders
	GetOrdersPagination(limit int, offset int, tags string) ([]*models.Order, int, error)
	GetStudentOrders(studentID string) ([]*models.Order, error)
	GetStudentOrdersPagination(limit int, offset int, studentID string) ([]*models.Order, int, error)

	// Response Edit
	CreateResponse(orderID string, response *models.NewResponseModel, Tutor *models.Tutor, username string) (string, error)

	// Response Getting
	GetResponsesByOrderID(id string) ([]models.Response, error)
	GetTutorsResponses(tutorID string) ([]models.Response, error)

	// Helpers
	GetTutorIsRespond(orderID string, tutorID string) (bool, error)
	GetUserByOrder(orderID string) (string, error)
	CheckOrderByStudentID(orderID string, studentID string) (bool, error)
	CheckResponseExist(TutorID, OrderID string) bool
	GetResponseById(ResponseID string) (*models.ResponseDB, error)

	GetDB() *sqlx.DB
}
