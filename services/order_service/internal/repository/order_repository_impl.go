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
	"github.com/randnull/Lessons/pkg/custom_errors"
	custom_logger "github.com/randnull/Lessons/pkg/logger"
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

	log.Printf("Connecting to database: %s", DbDatabase)

	db, err := sqlx.Open("postgres", link)

	if err != nil {
		log.Fatal("[Postgres] failed to connect" + err.Error())
	}

	PingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(PingCtx)

	if err != nil {
		log.Fatal("[Postgres] failed to ping" + err.Error())
	}

	log.Print("[Postgres] Database is ready")

	return &Repository{
		db: db,
	}
}

func MapDBError(err error, funcName string) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return custom_errors.ErrorNotFound
	}

	custom_logger.Info(fmt.Sprintf("[Postgres] %v failed: %v", funcName, err.Error()))

	return custom_errors.ErrorServiceError
}

func (o *Repository) CreateOrder(NewOrder *models.CreateOrder) (string, error) {
	const query = `
		INSERT INTO orders
			(
			 name, 
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
			 )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)  RETURNING id`

	var orderID string
	tags := pq.Array(NewOrder.Order.Tags)
	timestamp := time.Now()

	err := o.db.QueryRow(query,
		NewOrder.Order.Name,
		NewOrder.StudentID,
		NewOrder.Order.Title,
		NewOrder.Order.Description,
		NewOrder.Order.Grade,
		tags,
		NewOrder.Order.MinPrice,
		NewOrder.Order.MaxPrice,
		models.StatusWaiting,
		0,
		timestamp,
		timestamp,
	).Scan(&orderID)

	if err != nil {
		return "", MapDBError(err, "CreateOrder")
	}

	return orderID, nil
}

func (o *Repository) UpdateOrder(orderID string, order *models.UpdateOrder) error {
	query := `UPDATE orders SET `

	var values []interface{}

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

	currentTime := time.Now()

	query += fmt.Sprintf(`updated_at = $%v, `, strconv.Itoa(index))
	values = append(values, currentTime)
	index += 1

	if index == 2 {
		return nil
	}

	query = query[:len(query)-2] + ` WHERE id = $` + strconv.Itoa(index)
	values = append(values, orderID)

	_, err := o.db.Exec(query, values...)

	if err != nil {
		return MapDBError(err, "UpdateOrder")
	}

	return nil
}

func (o *Repository) DeleteOrder(id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	_, err := o.db.Exec(query, id)

	if err != nil {
		return MapDBError(err, "DeleteOrder")
	}

	return nil
}

func (o *Repository) SetOrderStatus(status string, orderID string) error {
	querySetStatus := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := o.db.Exec(querySetStatus, status, orderID)

	if err != nil {
		return MapDBError(err, "SetOrderStatus")
	}

	return nil
}

func (o *Repository) SetTutorToOrder(response *models.ResponseDB, UserData models.UserData) error {
	queryCheckStatus := `SELECT status FROM orders WHERE id = $1`

	var status string

	err := o.db.QueryRow(queryCheckStatus, response.OrderID).Scan(&status)

	if status != models.StatusNew {
		return custom_errors.ErrorAlreadySetTutor
	}

	tx, err := o.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return MapDBError(err, "SetTutorToOrder")
	}

	querySetStatus := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err = tx.Exec(querySetStatus, models.StatusSelected, response.OrderID)

	if err != nil {
		return MapDBError(err, "SetTutorToOrder")
	}

	queryUpdateResponses := `UPDATE responses SET is_final = $1 WHERE id = $2`

	_, err = tx.Exec(queryUpdateResponses, true, response.ID)

	if err != nil {
		return MapDBError(err, "SetTutorToOrder")
	}

	err = tx.Commit()

	if err != nil {
		return MapDBError(err, "SetTutorToOrder")
	}

	return nil
}

// Order getting

func (o *Repository) GetOrderByID(id string) (*models.Order, error) {
	const query = `
		SELECT 
			id, 
			name,
			student_id,
			title, 
			description, 
			grade,
			min_price, 
			max_price, 
			tags,
			status,
			response_count,
			created_at,
			updated_at
		FROM orders WHERE id = $1`

	var order models.Order

	err := o.db.QueryRowx(query, id).StructScan(&order)

	if err != nil {
		return nil, MapDBError(err, "GetOrderByID")
	}

	return &order, nil
}

func (o *Repository) GetOrders() ([]*models.Order, error) {
	var orders []*models.Order

	query := `
		SELECT 
			id, 
			name,
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

	err := o.db.Select(&orders, query)

	if err != nil {
		return nil, MapDBError(err, "GetOrders")
	}

	return orders, nil
}

func (o *Repository) GetOrdersPagination(limit int, offset int, tags string) ([]*models.Order, int, error) {
	var orders []*models.Order

	var args []interface{}
	var countArgs []interface{}

	countQuery := `
		SELECT
		    COUNT(*)
		FROM orders WHERE status = $1`

	countArgs = append(countArgs, models.StatusNew)

	if tags != "" {
		countQuery += ` AND $2 = ANY(tags)`
		countArgs = append(countArgs, tags)
	}

	var total int

	err := o.db.Get(&total, countQuery, countArgs...)

	if err != nil {
		return nil, 0, MapDBError(err, "GetOrdersPagination")
	}

	query := `
		SELECT 
			id, 
			name,
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
		FROM orders 
		WHERE status = $1`

	args = append(args, models.StatusNew)

	if tags != "" {
		query += ` AND $2 = ANY(tags)`
		args = append(args, tags)
		query += ` ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	} else {
		query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	}

	args = append(args, limit, offset)

	err = o.db.Select(&orders, query, args...)

	if err != nil {
		return nil, 0, MapDBError(err, "GetOrdersPagination")
	}

	return orders, total, nil
}

func (o *Repository) GetStudentOrdersPagination(limit int, offset int, studentID string) ([]*models.Order, int, error) {
	const queryCountOrders = `
		SELECT 
		    COUNT(*) 
		FROM orders 
		WHERE student_id = $1`

	var total int

	err := o.db.Get(&total, queryCountOrders, studentID)

	if err != nil {
		return nil, 0, MapDBError(err, "GetStudentOrdersPagination")
	}

	var orders []*models.Order

	const querySelectOrders = `
		SELECT 
			id, 
			name,
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
		FROM orders 
		WHERE student_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	err = o.db.Select(&orders, querySelectOrders, studentID, limit, offset)

	if err != nil {
		return nil, 0, MapDBError(err, "GetStudentOrdersPagination")
	}

	return orders, total, nil
}

func (o *Repository) GetStudentOrders(studentID string) ([]*models.Order, error) {
	var orders []*models.Order

	query := `
		SELECT 
    		id, 
    		name,
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

	err := o.db.Select(&orders, query, studentID)

	if err != nil {
		return nil, MapDBError(err, "GetStudentOrders")
	}

	return orders, nil
}

func (o *Repository) CreateResponse(orderID string,
	response *models.NewResponseModel,
	Tutor *models.Tutor,
	username string) (string, error) {
	var ResponseID string

	const queryCheck = `
		SELECT
    		id 
		FROM responses 
		WHERE order_id = $1 AND tutor_id = $2`

	err := o.db.QueryRow(queryCheck, orderID, Tutor.Id).Scan(&ResponseID)

	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			custom_logger.Info("[Postgres] CreateResponse: ErrResponseAlreadyExist")
			return ResponseID, custom_errors.ErrResponseAlredyExist
		}
		custom_logger.Error("[Postgres] CreateResponse error" + err.Error())
		return "", err
	}

	tx, err := o.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return "", MapDBError(err, "CreateResponse")
	}

	timestamp := time.Now()

	const queryInsert = `
		INSERT INTO responses 
		    (
		     order_id, 
		     name, 
		     tutor_id, 
		     tutor_username, 
		     greetings, 
		     is_final, 
		     created_at
		     )
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	greetingsMessage := response.Greetings

	err = tx.QueryRow(queryInsert, orderID, Tutor.Name, Tutor.Id, username, greetingsMessage, false, timestamp).Scan(&ResponseID)

	if err != nil {
		return "", MapDBError(err, "CreateResponse")
	}

	const queryUpdate = `
		UPDATE orders SET
			response_count = response_count + 1
		WHERE id = $1`

	_, err = tx.Exec(queryUpdate, orderID)

	if err != nil {
		return "", MapDBError(err, "CreateResponse")
	}

	err = tx.Commit()

	if err != nil {
		return "", MapDBError(err, "CreateResponse")
	}

	return ResponseID, nil
}

func (o *Repository) GetResponsesByOrderID(id string) ([]models.Response, error) {
	const query = `
		SELECT 
			id,
			name,
			tutor_id,
			is_final,
			created_at,
			order_id
		FROM responses
		WHERE order_id = $1`

	var responses []models.Response

	err := o.db.Select(&responses, query, id)

	if err != nil {
		return nil, MapDBError(err, "GetResponsesByOrderID")
	}

	return responses, nil
}

func (o *Repository) GetTutorsResponses(tutorID string) ([]models.Response, error) {
	const query = `
		SELECT 
			id,
			name,
			tutor_id,
			is_final,
			created_at,
			order_id
		FROM responses
		WHERE tutor_id = $1`

	var responses []models.Response

	err := o.db.Select(&responses, query, tutorID)

	if err != nil {
		return nil, MapDBError(err, "GetTutorsResponses")
	}

	return responses, nil
}

func (o *Repository) GetResponseById(ResponseID string) (*models.ResponseDB, error) {
	const query = `
		SELECT 
			id,
			order_id,
			tutor_id,
			tutor_username,
			name,
			greetings,
			is_final,
			created_at
		FROM responses WHERE id = $1`

	var response models.ResponseDB

	err := o.db.QueryRowx(query, ResponseID).StructScan(&response)

	if err != nil {
		return nil, MapDBError(err, "GetResponseById")
	}

	return &response, nil
}

// Helpers

func (o *Repository) GetTutorIsRespond(orderID string, tutorID string) (bool, error) {
	const query = `
        SELECT EXISTS (
            SELECT 1 FROM responses WHERE order_id = $1 AND tutor_id = $2
        )`

	var isExist = false

	err := o.db.QueryRow(query, orderID, tutorID).Scan(&isExist)

	if err != nil {
		return false, MapDBError(err, "GetTutorIsRespond")
	}

	return isExist, nil
}

func (o *Repository) GetUserByOrder(orderID string) (string, error) {
	var UserID string

	query := `SELECT student_id FROM orders WHERE id = $1`

	err := o.db.QueryRow(query, orderID).Scan(&UserID)

	if err != nil {
		return "", MapDBError(err, "GetUserByOrder")
	}

	return UserID, nil
}

func (o *Repository) CheckResponseExist(tutorID, orderID string) bool {
	var isExist bool

	const query = `
		SELECT EXISTS (
			SELECT 1 FROM responses WHERE order_id = $1 AND tutor_id = $2
		)`

	err := o.db.QueryRow(query, orderID, tutorID).Scan(&isExist)

	if err != nil {
		return false
	}

	return isExist
}

func (o *Repository) CheckOrderByStudentID(orderID string, studentID string) (bool, error) {
	var isExist bool

	const query = `
        SELECT EXISTS (
            SELECT 1 FROM orders WHERE id = $1 AND student_id = $2
        )`

	err := o.db.QueryRow(query, orderID, studentID).Scan(&isExist)

	if err != nil {
		return false, MapDBError(err, "CheckOrderByStudentID")
	}

	return isExist, nil
}

func (o *Repository) GetDB() *sqlx.DB {
	return o.db
}
