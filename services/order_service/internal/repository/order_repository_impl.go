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
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/models"
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

func (orderStorage *Repository) CreateOrder(order *models.NewOrder, studentID string, telegramID int64) (*models.OrderToBrokerWithID, error) {
	timestamp := time.Now()

	query := `INSERT INTO orders (student_id, title, description, grade, tags, min_price, max_price, status, response_count, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)  RETURNING id`

	var orderID string

	// conver []string -> pq.stringArray

	tags := pq.Array(order.Tags)

	log.Println(tags)

	err := orderStorage.db.QueryRow(query,
		studentID,
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
		return nil, err
	} // Норм проверку TODO

	CreatedOrder := models.OrderToBrokerWithID{
		ID:          orderID,
		StudentID:   telegramID,
		Title:       order.Title,
		Description: order.Description,
		Tags:        order.Tags,
		MinPrice:    order.MinPrice,
		MaxPrice:    order.MaxPrice,
		ChatID:      telegramID,
	}

	return &CreatedOrder, nil
}

func (orderStorage *Repository) GetByID(id string, studentID string) (*models.OrderDetails, error) {
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
			r.name,
			r.tutor_id,
			r.is_final,
			r.created_at
		FROM orders o
		LEFT JOIN responses r ON o.id = r.order_id
		WHERE o.id = $1 AND o.student_id = $2`

	rows, err := orderStorage.db.Query(query, id, studentID)

	fmt.Println(rows, err)
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
		var responseID sql.NullString
		var tutorID sql.NullString
		var responseCreatedAt sql.NullTime
		var tutorName sql.NullString
		var isFinal sql.NullBool

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
			&tutorName,
			&tutorID,
			&isFinal,
			&responseCreatedAt,
		)

		if err != nil {
			fmt.Println(err)

			return nil, err
		}

		if responseID.Valid {
			validResponse := models.Response{
				ID:        responseID.String,
				Name:      tutorName.String,
				TutorID:   tutorID.String,
				IsFinal:   isFinal.Bool,
				CreatedAt: responseCreatedAt.Time,
			}
			responses = append(responses, validResponse)
		}
	}

	order.Responses = responses

	return order, nil
}

func (orderStorage *Repository) GetUserByOrder(orderID string) (string, error) {
	var UserID string

	query := `SELECT student_id FROM orders WHERE id = $1`

	err := orderStorage.db.QueryRow(query, orderID).Scan(&UserID)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return UserID, nil
}

func (orderStorage *Repository) GetAllOrders(studentID string) ([]*models.Order, error) {
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
			FROM orders WHERE status = $1 ORDER BY created_at DESC`

	rows, err := orderStorage.db.Query(query, "New")
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

func (orderStorage *Repository) GetOrderByIdTutor(id string, tutorID string) (*models.OrderDetailsTutor, error) {
	var order models.OrderDetailsTutor

	query := `
		SELECT 
			id, 
			title, 
			description, 
			grade,
			min_price, 
			max_price, 
			tags,
			status,
			response_count,
			created_at
		FROM orders WHERE id = $1`

	err := orderStorage.db.QueryRowx(query, id).StructScan(&order)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	query = `
        SELECT EXISTS (
            SELECT 1 FROM responses WHERE order_id = $1 AND tutor_id = $2
        )
    `

	var isExist = false

	err = orderStorage.db.QueryRow(query, id, tutorID).Scan(&isExist)

	if err != nil {
		isExist = false
	}

	order.IsResponsed = isExist

	if err != nil {
		return nil, custom_errors.ErrGetOrder
	}

	return &order, nil
}

func (orderStorage *Repository) GetAllUsersOrders(studentID string) ([]*models.Order, error) {
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

	rows, err := orderStorage.db.Query(query, studentID)
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return orders, nil
}

func (orderStorage *Repository) UpdateOrder(orderID string, order *models.UpdateOrder, studentID string) error {
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

func (orderStorage *Repository) DeleteOrder(id string, studentID string) error {
	query := `DELETE FROM orders WHERE id = $1 AND student_id = $2`

	_, err := orderStorage.db.Exec(query, id, studentID)

	if err != nil {
		return err
	}

	return nil
}

//func (orderStorage *Repository) VerifyUserOrder(studentID string) ([]*models.Order, error) {}

func (orderStorage *Repository) CreateResponse(response *models.NewResponseModel, Tutor *models.User, username string) (string, error) {
	var ResponseID string

	queryCheck := `SELECT id FROM responses WHERE order_id = $1 AND tutor_id = $2`

	err := orderStorage.db.QueryRow(queryCheck, response.OrderId, Tutor.Id).Scan(&ResponseID)

	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return ResponseID, custom_errors.ErrResponseAlredyExist
		}
		return "", err
	}

	tx, err := orderStorage.db.Begin()
	defer tx.Rollback()

	if err != nil {
		tx.Rollback()
		return "", err
	}

	timestamp := time.Now()
	// SELECT WHERE order_id = ... без UPDATE
	queryInsert := `INSERT INTO responses (order_id, name, tutor_id, tutor_username, is_final, created_at)
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = tx.QueryRow(queryInsert, response.OrderId, Tutor.Name, Tutor.Id, username, false, timestamp).Scan(&ResponseID)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return "", err
	}

	// response_count = response_count + 1
	queryUpdate := `UPDATE orders SET response_count = response_count + 1 WHERE id = $1`

	_, err = tx.Exec(queryUpdate, response.OrderId)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return ResponseID, nil
}

func (orderStorage *Repository) CheckOrderByStudentID(orderID string, studentID string) (bool, error) {
	var isExist bool

	query := `
        SELECT EXISTS (
            SELECT 1 FROM orders WHERE id = $1 AND student_id = $2
        )
    `

	err := orderStorage.db.QueryRow(query, orderID, studentID).Scan(&isExist)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (orderStorage *Repository) GetResponseById(ResponseID string, studentID string) (*models.ResponseDB, error) {
	var response models.ResponseDB

	query := `
		SELECT 
			id,
			order_id,
			tutor_id,
			tutor_username,
			name,
			is_final,
			created_at
		FROM responses WHERE id = $1`

	err := orderStorage.db.QueryRowx(query, ResponseID).StructScan(&response)

	if err != nil {
		fmt.Println(err)
		return nil, custom_errors.ErrGetResponse
	}

	fmt.Println(response.OrderID, studentID)

	isExist, err := orderStorage.CheckOrderByStudentID(response.OrderID, studentID)

	if err != nil {
		return nil, err
	}

	if !isExist {
		return nil, custom_errors.ErrNotAllowed
	}

	return &response, nil
}

func (orderStorage *Repository) SetTutorToOrder(responseID string, UserData models.UserData) error {
	response, err := orderStorage.GetResponseById(responseID, UserData.UserID)

	log.Println(response, err)

	if err != nil || response == nil {
		return custom_errors.ErrNotAllowed
	}

	isExist, err := orderStorage.CheckOrderByStudentID(response.OrderID, UserData.UserID)

	log.Println(isExist, err)

	if err != nil || !isExist {
		return custom_errors.ErrNotAllowed
	}

	queryCheckStatus := `SELECT status FROM orders WHERE id = $1`

	log.Println(queryCheckStatus)

	var status string

	err = orderStorage.db.QueryRow(queryCheckStatus, response.OrderID).Scan(&status)

	log.Println(err, status)

	if status != "New" {
		return custom_errors.ErrorAlreadySetTutor
	}

	tx, err := orderStorage.db.Begin()
	defer tx.Rollback()
	log.Println(err)

	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("Selected", response.ID, response.OrderID)

	querySetStatus := `UPDATE orders SET status = $1 WHERE id = $2` // status = $1,

	_, err = tx.Exec(querySetStatus, "Selected", response.OrderID) // "Selected",

	log.Println(err)
	if err != nil {
		tx.Rollback()
		return err
	}

	queryUpdateResponses := `UPDATE responses SET is_final = $1 WHERE id = $2`

	_, err = tx.Exec(queryUpdateResponses, true, response.ID)
	log.Println(err)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	log.Println(err)

	if err != nil {
		return err
	}

	return nil
}
