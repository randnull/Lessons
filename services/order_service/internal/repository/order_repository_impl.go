package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/randnull/Lessons/internal/models"
	"log"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository() *Repository {

	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"CHANGE", "CHANGE", "postgresql", "5432", "orders_database")

	db, err := sqlx.Open("postgres", link)

	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Database is ready")

	return &Repository{
		db: db,
	}
}

func (orderStorage *Repository) CreateOrder(order *models.NewOrder) (string, error) {
	timestamp := time.Now()

	fmt.Println(timestamp)

	query := `INSERT INTO orders (student_id, tutor_id, subject, description, min_price, max_price, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)  RETURNING id`

	var orderID string

	err := orderStorage.db.QueryRow(query, "61455107-52cc-4fe4-9c95-3110e65d1a06", "61455107-52cc-4fe4-9c95-3110e65d1a06", order.Subject, order.Description, order.MinPrice, order.MaxPrice, timestamp, timestamp).Scan(&orderID)

	if err != nil {
		log.Fatal(err)
	}

	return orderID, err
}

func (orderStorage *Repository) GetByID(id string) (*models.Order, error) {
	order := &models.Order{}

	query := `SELECT id, student_id, tutor_id, subject, description, min_price, max_price, created_at, updated_at FROM orders WHERE id = $1`

	err := orderStorage.db.QueryRow(query, id).Scan(&order.ID, &order.StudentID, &order.TutorID, &order.Subject, &order.Description, &order.MinPrice, &order.MaxPrice, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (orderStorage *Repository) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT id, student_id, tutor_id, subject, description, min_price, max_price, created_at, updated_at FROM orders`

	rows, err := orderStorage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.StudentID, &order.TutorID, &order.Subject, &order.Description, &order.MinPrice, &order.MaxPrice, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
