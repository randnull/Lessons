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
	DbUser := cfg.DBUser
	DbPassword := cfg.DBPassword
	DbHost := cfg.DBHost
	DbPort := cfg.DBPort
	DbDatabase := cfg.DBName

	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		DbUser,
		DbPassword,
		DbHost,
		DbPort,
		DbDatabase,
	)

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

func (orderStorage *Repository) CreateOrder(order *models.NewOrder, InitData initdata.InitData) (models.OrderToBrokerWithID, error) {
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
		"New", // этот кринж на enum TODO
		timestamp,
		timestamp,
	).Scan(&orderID)

	if err != nil {
		log.Fatal(err)
	} // Норм проверку TODO

	CreatedOrder := models.OrderToBrokerWithID{
		ID:          orderID,
		StudentID:   int(InitData.User.ID),
		Title:       order.Title,
		Description: order.Description,
		Tags:        order.Tags,
		MinPrice:    order.MinPrice,
		MaxPrice:    order.MaxPrice,
		ChatID:      InitData.Chat.ID,
	}

	return CreatedOrder, nil
}

func (orderStorage *Repository) GetByID(id string, InitData initdata.InitData) (*models.Order, error) {
	order := &models.Order{}

	query := `SELECT 
    			id, 
    			student_id, 
    			title, 
    			description, 
    			tags, 
    			min_price, 
    			max_price, 
    			status, 
    			created_at, 
    			updated_at 
			FROM orders WHERE id = $1 AND student_id = $2`

	err := orderStorage.db.QueryRow(query, id, InitData.User.ID).Scan(
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
			return nil, errors.New("Order not found")
		} else {
			return nil, err
		}
	}

	return order, nil
}

func (orderStorage *Repository) GetAllOrders(InitData initdata.InitData) ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT 
    			id, 
    			student_id, 
    			title, 
    			description, 
    			tags, 
    			min_price, 
    			max_price, 
    			status, 
    			created_at, 
    			updated_at 
			FROM orders WHERE student_id = $1`

	rows, err := orderStorage.db.Query(query, InitData.User.ID)
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

	for i, j := 0, len(orders)-1; i < j; i, j = i+1, j-1 {
		orders[i], orders[j] = orders[j], orders[i]
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (orderStorage *Repository) UpdateOrder(orderID string, order *models.NewOrder, InitData initdata.InitData) error {
	return nil
}

func (orderStorage *Repository) DeleteOrder(id string, InitData initdata.InitData) error {
	query := `DELETE FROM orders WHERE id = $1 AND student_id = $2`

	_, err := orderStorage.db.Exec(query, id, InitData.User.ID)

	if err != nil {
		return err
	}

	return nil
}

//func (orderStorage *Repository) VerifyUserOrder(studentID string) ([]*models.Order, error) {}
