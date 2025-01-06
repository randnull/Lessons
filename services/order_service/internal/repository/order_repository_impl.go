package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"log"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(cfg config.DBConfig) *Repository {

	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		"change", "9yuVZktnLKzqMrkywVgTlhDxVQsqWXbP", "dpg-cttubetumphs73eikdbg-a.oregon-postgres.render.com", "5432", "orders_database_bhw2")

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

func (orderStorage *Repository) CreateOrder(order *models.NewOrder, InitData initdata.InitData) (string, error) {
	timestamp := time.Now()

	fmt.Println(timestamp)
	//student_id, tutor_id, subject, description, min_price, max_price, created_at, updated_at
	query := `INSERT INTO orders (student_id, title, description, tags, min_price, max_price, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)  RETURNING id`

	var orderID string

	// conver []string -> pq.stringArray

	tags := pq.Array(order.Tags)

	log.Println(tags)

	err := orderStorage.db.QueryRow(query,
		InitData.User.ID,
		order.Title,
		order.Description,
		tags,
		order.MinPrice,
		order.MaxPrice,
		"New",
		timestamp,
		timestamp,
	).Scan(&orderID)

	if err != nil {
		log.Fatal(err)
	} // Норм проверку

	return orderID, nil
}

func (orderStorage *Repository) GetByID(id string) (*models.Order, error) {
	order := &models.Order{}

	query := `SELECT id, student_id, title, description, tags, min_price, max_price, status, created_at, updated_at FROM orders WHERE id = $1`

	err := orderStorage.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.StudentID,
		&order.Title,
		&order.Description,
		&order.Tags,
		&order.MinPrice,
		&order.MaxPrice,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Order not found!")
		} else {
			log.Fatal(err)
		}
	}

	return order, nil
}

func (orderStorage *Repository) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT id, student_id, title, description, tags, min_price, max_price, status, created_at, updated_at FROM orders`

	rows, err := orderStorage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.StudentID,
			&order.Title,
			&order.Description,
			&order.Tags,
			&order.MinPrice,
			&order.MaxPrice,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
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

func (orderStorage *Repository) UpdateOrder(order *models.Order) error {
	return nil
}

func (orderStorage *Repository) DeleteOrder(id string) error {
	return nil
}
