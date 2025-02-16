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
	"strconv"
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

	link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
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

	query := `INSERT INTO orders (student_id, title, description, grade, tags, min_price, max_price, status, response_count, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)  RETURNING id`

	var orderID string

	// conver []string -> pq.stringArray

	tags := pq.Array(order.Tags)

	log.Println(tags)

	err := orderStorage.db.QueryRow(query,
		InitData.User.ID,
		order.Title,
		order.Description,
		order.Grade,
		tags,
		order.MinPrice,
		order.MaxPrice,
		"New", // этот кринж на enum TODO
		0,
		timestamp,
		timestamp,
	).Scan(&orderID)

	if err != nil {
		log.Println(err)
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

func (orderStorage *Repository) GetByID(id string, InitData initdata.InitData) (*models.OrderDetails, error) {
	order := &models.OrderDetails{}
	responses := []models.Response{}

	query := `
		SELECT 
			o.id, 
			o.student_id, 
			o.title, 
			o.description, 
			o.grade,
			o.tags, 
			o.min_price, 
			o.max_price, 
			o.status,
			o.response_count,
			o.created_at, 
			o.updated_at,
			r.id,
			r.tutor_id,
			r.created_at
		FROM orders o
		LEFT JOIN responses r ON o.id = r.order_id
		WHERE o.id = $1 AND o.student_id = $2`

	rows, err := orderStorage.db.Query(query, id, InitData.User.ID)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Order not found!")
			return nil, errors.New("Order not found")
		} else {
			return nil, err
		}
	}
	for rows.Next() {
		var responseID sql.NullString //переписать !!!
		var tutorID sql.NullInt64
		var responseCreatedAt sql.NullTime
		err := rows.Scan(
			&order.ID,
			&order.StudentID,
			&order.Title,
			&order.Description,
			&order.Grade,
			&order.Tags,
			&order.MinPrice,
			&order.MaxPrice,
			&order.Status,
			&order.ResponseCount,
			&order.CreatedAt,
			&order.UpdatedAt,
			&responseID,
			&tutorID,
			&responseCreatedAt,
		)

		fmt.Println(responseID, tutorID, responseCreatedAt)

		if err != nil {
			fmt.Println(err)

			return nil, err
		}

		if responseID.Valid {
			valid_response := models.Response{
				ID:        responseID.String,
				TutorID:   tutorID.Int64,
				CreatedAt: responseCreatedAt.Time,
			}
			responses = append(responses, valid_response)
		}
	}
	// Request to userservice
	for index := range responses {
		responses[index].Name = "John"
	}

	order.Responses = responses

	fmt.Println(order)

	fmt.Println(err)

	return order, nil
}

func (orderStorage *Repository) GetAllOrders(InitData initdata.InitData) ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT 
    			id, 
    			student_id, 
    			title, 
    			description,
    			grade,
    			tags, 
    			min_price, 
    			max_price, 
    			status,
    			response_count,
    			created_at, 
    			updated_at 
			FROM orders ORDER BY created_at DESC`

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
			&order.Grade,
			&order.Tags,
			&order.MinPrice,
			&order.MaxPrice,
			&order.Status,
			&order.ResponseCount,
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

func (orderStorage *Repository) GetAllUsersOrders(InitData initdata.InitData) ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT 
    			id, 
    			student_id, 
    			title, 
    			description, 
    			grade,
    			tags, 
    			min_price, 
    			max_price, 
    			status, 
    			response_count,
    			created_at, 
    			updated_at 
			FROM orders WHERE student_id = $1 ORDER BY created_at DESC`

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
			&order.Grade,
			&order.Tags,
			&order.MinPrice,
			&order.MaxPrice,
			&order.Status,
			&order.ResponseCount,
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

func (orderStorage *Repository) UpdateOrder(orderID string, order *models.UpdateOrder, InitData initdata.InitData) error {
	query := `UPDATE orders SET `
	values := []interface{}{}

	index := 1

	if order.Title != "" {
		query += fmt.Sprintf(`title = $%v, `, strconv.Itoa(index))
		values = append(values, order.Title)
		index += 1
	}

	if order.Description != "" {
		query += fmt.Sprintf(`description = $%v, `, strconv.Itoa(index))
		values = append(values, order.Description)
		index += 1
	}

	if order.Grade != "" {
		query += fmt.Sprintf(`grade = $%v, `, strconv.Itoa(index))
		values = append(values, order.Grade)
		index += 1
	}

	//if order.MinPrice != 0 {
	//	query += fmt.Sprintf(`min_price = $%v, `, strconv.Itoa(index))
	//	values = append(values, order.MinPrice)
	//	index += 1
	//}
	//
	//if order.MaxPrice != 0 {
	//	query += fmt.Sprintf(`max_price = $%v, `, strconv.Itoa(index))
	//	values = append(values, order.MaxPrice)
	//	index += 1
	//}

	if index == 1 {
		return nil
	}

	query = query[:len(query)-2] + ` WHERE id = $` + strconv.Itoa(index)
	values = append(values, orderID)

	_, err := orderStorage.db.Exec(query, values...)

	if err != nil {
		return err
	}

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

func (orderStorage *Repository) CreateResponse(response *models.NewResponseModel, InitData initdata.InitData) error {
	tx, err := orderStorage.db.Begin()
	if err != nil {
		return err
	}

	fmt.Println("response start")

	timestamp := time.Now()
	// SELECT WHERE order_id = ... без UPDATE
	queryInsert := `INSERT INTO responses (order_id, tutor_id, created_at) VALUES ($1, $2, $3) RETURNING id`

	//_, err := orderStorage.db.Exec(query, response.OrderId, response.OrderId, timestamp)
	_, err = tx.Exec(queryInsert, response.OrderId, InitData.User.ID, timestamp)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	fmt.Println("response start 2")

	// response_count = response_count + 1
	queryUpdate := `UPDATE orders SET response_count = response_count + 1 WHERE id = $1`
	fmt.Println(response.OrderId)
	res, err := tx.Exec(queryUpdate, response.OrderId)
	fmt.Println(res, err)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("response end")

	return nil
}
