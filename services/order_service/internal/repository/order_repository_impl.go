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

	query := `INSERT INTO orders (
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

	tags := pq.Array(order.Tags)

	log.Println(tags)

	err := orderStorage.db.QueryRow(query,
		order.Name,
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
		ID:        orderID,
		StudentID: telegramID,
		Title:     order.Title,
		Tags:      order.Tags,
		Status:    "New",
	}

	return &CreatedOrder, nil
}

func (orderStorage *Repository) GetOrderByID(id string) (*models.OrderDetails, error) {
	order := &models.OrderDetails{}
	responses := []models.Response{}

	query := `
		SELECT 
			o.id, 
			o.name,
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
		WHERE o.id = $1`

	rows, err := orderStorage.db.Query(query, id)

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
			&order.Name,
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
				OrderID:   id,
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

func (orderStorage *Repository) GetOrders() ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT 
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
			&order.Name,
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

func (orderStorage *Repository) GetOrdersPagination(limit int, offset int, tags string) ([]*models.Order, int, error) {
	tx, err := orderStorage.db.Begin()

	if err != nil {
		return nil, 0, err
	}

	defer tx.Rollback()

	var total int

	queryCount := `SELECT 
    					COUNT(*) 
					FROM orders WHERE status = $1`

	countArgs := []interface{}{"New"}

	if tags != "" {
		queryCount += ` AND $2 = ANY(tags)`
		countArgs = append(countArgs, tags)
	}

	err = tx.QueryRow(queryCount, countArgs...).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	var orders []*models.Order

	query := `SELECT 
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
            FROM orders WHERE status = $1`

	queryArgs := []interface{}{"New"}

	if tags != "" {
		query += ` AND $2 = ANY(tags)`
		queryArgs = append(queryArgs, tags)
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(queryArgs)+1) + ` OFFSET $` + strconv.Itoa(len(queryArgs)+2)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := tx.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.Name,
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
			return nil, 0, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	if err = tx.Commit(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (orderStorage *Repository) GetStudentOrdersPagination(limit int, offset int, studentID string) ([]*models.Order, int, error) {
	tx, err := orderStorage.db.Begin()

	if err != nil {
		return nil, 0, err
	}

	defer tx.Rollback()

	var total int

	queryCount := `SELECT 
    					COUNT(*) 
					FROM orders WHERE student_id = $1`

	err = tx.QueryRow(queryCount, studentID).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	var orders []*models.Order

	query := `SELECT 
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
			FROM orders WHERE student_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := tx.Query(query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.Name,
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
			return nil, 0, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	if err = tx.Commit(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (orderStorage *Repository) GetOrderByIdTutor(id string, tutorID string) (*models.OrderDetailsTutor, error) {
	var order models.OrderDetailsTutor

	query := `
		SELECT 
			id, 
			name,
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

func (orderStorage *Repository) GetStudentOrders(studentID string) ([]*models.Order, error) {
	var orders []*models.Order

	query := `SELECT 
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
			&order.Name,
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

func (orderStorage *Repository) DeleteOrder(id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	_, err := orderStorage.db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (orderStorage *Repository) CheckResponseExist(TutorID, OrderID string) bool {
	var ResponseID string

	queryCheck := `SELECT id FROM responses WHERE order_id = $1 AND tutor_id = $2`

	err := orderStorage.db.QueryRow(queryCheck, OrderID, TutorID).Scan(&ResponseID)

	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		if err == nil {
			return true
		}
	}
	return false
}

func (orderStorage *Repository) CreateResponse(orderID string,
	response *models.NewResponseModel,
	Tutor *models.Tutor,
	username string) (string, error) {
	var ResponseID string

	queryCheck := `SELECT id FROM responses WHERE order_id = $1 AND tutor_id = $2`

	err := orderStorage.db.QueryRow(queryCheck, orderID, Tutor.Id).Scan(&ResponseID)

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

	queryInsert := `INSERT INTO responses (order_id, name, tutor_id, tutor_username, greetings, is_final, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	greetingsMessage := response.Greetings

	err = tx.QueryRow(queryInsert, orderID, Tutor.Name, Tutor.Id, username, greetingsMessage, false, timestamp).Scan(&ResponseID)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return "", err
	}

	queryUpdate := `UPDATE orders SET response_count = response_count + 1 WHERE id = $1`

	_, err = tx.Exec(queryUpdate, orderID)

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

func (orderStorage *Repository) GetResponseById(ResponseID string) (*models.ResponseDB, error) {
	var response models.ResponseDB

	query := `
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

	err := orderStorage.db.QueryRowx(query, ResponseID).StructScan(&response)

	if err != nil {
		fmt.Println(err)
		return nil, custom_errors.ErrGetResponse
	}

	return &response, nil
}

func (orderStorage *Repository) SetTutorToOrder(response *models.ResponseDB, UserData models.UserData) error {
	queryCheckStatus := `SELECT status FROM orders WHERE id = $1`

	var status string

	err := orderStorage.db.QueryRow(queryCheckStatus, response.OrderID).Scan(&status)

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

	querySetStatus := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err = tx.Exec(querySetStatus, "Selected", response.OrderID) // "Selected",

	log.Println(err)
	if err != nil {
		tx.Rollback()
		return err
	}

	queryUpdateResponses := `UPDATE responses SET is_final = $1 WHERE id = $2`

	_, err = tx.Exec(queryUpdateResponses, true, response.ID)

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

func (orderStorage *Repository) SetOrderStatus(status string, orderID string) error {
	querySetStatus := `UPDATE orders SET status = $1 WHERE id = $2`

	_, err := orderStorage.db.Exec(querySetStatus, status, orderID)

	if err != nil {
		log.Println(err)
		return errors.New("cannot set status")
	}
	return nil
}

func (orderStorage *Repository) GetTutorsResponses(tutorID string) ([]models.Response, error) {
	query := `
        SELECT 
            id,
            order_id,
            name,
            tutor_id,
            is_final,
            created_at
        FROM responses 
        WHERE tutor_id = $1
        ORDER BY created_at DESC`

	rows, err := orderStorage.db.Query(query, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.Response
	for rows.Next() {
		var resp models.Response
		err := rows.Scan(
			&resp.ID,
			&resp.OrderID,
			&resp.Name,
			&resp.TutorID,
			&resp.IsFinal,
			&resp.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return responses, nil
}
