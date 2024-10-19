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

func (orderStorage *Repository) CreateOrder(order *models.NewOrder) error {
	timestamp := time.Now()

	fmt.Println(timestamp)

	query := `INSERT INTO orders (student_id, tutor_id, subject, description, min_price, max_price, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := orderStorage.db.Exec(query, order.StudentID, order.TutorID, order.Subject, order.Description, order.MinPrice, order.MaxPrice, timestamp, timestamp)

	if err != nil {
		log.Fatal(err)
	}
	return err
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
