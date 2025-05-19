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
	pb "github.com/randnull/Lessons/pkg/gRPC"
	custom_logger "github.com/randnull/Lessons/pkg/logger"
	"strings"
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
		custom_logger.Error("[Postgres] failed to connect" + err.Error())
		return nil
	}

	err = db.PingContext(context.Background())

	if err != nil {
		custom_logger.Error("[Postgres] failed to ping" + err.Error())
		return nil
	}

	custom_logger.Info("[Postgres] Database is ready")

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

func (r *Repository) CreateUser(user *models.CreateUser) (string, error) {
	custom_logger.Info(fmt.Sprintf("[Postgres] CreateUser called. UserID: %v. Role: %v", user.TelegramId, user.Role))

	ExistedUser, err := r.GetUserByTelegramId(user.TelegramId, user.Role)

	if err == nil {
		custom_logger.Info(fmt.Sprintf("[Postgres] CreateUser User exist. Authorization Telegram-UserID: %v. UserID: %v", user.TelegramId, ExistedUser.Id))
		return ExistedUser.Id, nil
	}

	tx, err := r.db.Begin()
	defer tx.Rollback()

	if err != nil {
		return "", err
	}

	const query = `
		INSERT INTO users (
            telegram_id,
            name,
            role,
            created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var UserId string

	currentTime := time.Now()

	err = tx.QueryRow(query,
		user.TelegramId,
		user.Name,
		user.Role,
		currentTime,
	).Scan(&UserId)

	if err != nil {
		return "", MapDBError(err, "CreateUser")
	}

	if user.Role == models.RoleTutor {
		const queryInsertTutor = `
				INSERT INTO tutors 
				(
                    id,
                    created_at
                )
				VALUES ($1, $2)`

		_, err = tx.Exec(queryInsertTutor, UserId, currentTime)

		if err != nil {
			custom_logger.Info("[Postgres] CreateUser failed: " + err.Error())
			return "", MapDBError(err, "CreateUser")
		}
	}

	err = tx.Commit()

	if err != nil {
		return "", MapDBError(err, "CreateUser")
	}

	return UserId, nil
}

func (r *Repository) GetUserById(userID string) (*models.UserDB, error) {
	custom_logger.Info("[Postgres] GetUserById called. UserID: " + userID)

	var user models.UserDB

	const query = `
		SELECT 
		    id,
		    telegram_id,
		    name, 
		    role, 
		    is_banned,
		    created_at 
		FROM users 
		WHERE id = $1`

	err := r.db.Get(&user, query, userID)

	if err != nil {
		return nil, MapDBError(err, "GetUserById")
	}

	custom_logger.Info("[Postgres] GetUserById success. UserID: " + userID)

	return &user, nil
}

func (r *Repository) GetUserByTelegramId(telegramID int64, userRole string) (*models.UserDB, error) {
	var user models.UserDB

	const query = `
		SELECT 
		    id,
		    telegram_id,
		    name, 
		    role, 
		    is_banned,
		    created_at 
		FROM users 
		WHERE telegram_id = $1 AND role = $2`

	err := r.db.Get(&user, query, telegramID, userRole)

	if err != nil {
		return nil, MapDBError(err, "GetUserByTelegramId")
	}

	return &user, nil
}

func (r *Repository) GetStudentById(userID string) (*models.UserDB, error) {
	var user models.UserDB

	const query = `
		SELECT 
		    id,
		    telegram_id,
		    name, 
		    role, 
		    is_banned,
		    created_at 
		FROM users 
		WHERE id = $1 AND role = $2`

	err := r.db.Get(&user, query, userID, models.RoleStudent)

	if err != nil {
		return nil, MapDBError(err, "GetStudentById")
	}

	return &user, nil
}

func (r *Repository) GetTutorByID(userID string) (*models.TutorDB, error) {
	const query = `
        SELECT 
            u.id, 
            u.telegram_id, 
            u.name, 
            u.role, 
            u.created_at,
            t.bio,
            t.response_count,
            t.rating,
            t.tags,
            t.is_active,
            t.created_at
        FROM users u 
        INNER JOIN tutors t ON u.id = t.id
        WHERE u.id = $1 AND u.role = $2`

	var tutor models.TutorDB

	var bio sql.NullString
	var responseCount sql.NullInt32
	var tags pq.StringArray
	var isActive sql.NullBool
	var tutorCreatedAt sql.NullTime
	var rating sql.NullInt32

	err := r.db.QueryRow(query, userID, models.RoleTutor).Scan(
		&tutor.Id,
		&tutor.TelegramID,
		&tutor.Name,
		&tutor.Role,
		&tutor.CreatedAt,
		&bio,
		&responseCount,
		&rating,
		&tags,
		&isActive,
		&tutorCreatedAt,
	)

	if err != nil {
		return nil, MapDBError(err, "GetTutorByID")
	}

	if bio.Valid {
		tutor.Bio = bio.String
	}

	if responseCount.Valid {
		tutor.ResponseCount = responseCount.Int32
	}

	if tags != nil {
		tutor.Tags = tags
	} else {
		tutor.Tags = []string{}
	}

	if isActive.Valid {
		tutor.IsActive = isActive.Bool
	}

	if tutorCreatedAt.Valid {
		tutor.TutorCreatedAt = tutorCreatedAt.Time
	}

	if rating.Valid {
		tutor.Rating = rating.Int32
	}

	return &tutor, nil
}

func (r *Repository) GetAllUsers() ([]*pb.User, error) {
	const query = `
		SELECT 
			id, 
			name, 
			telegram_id, 
			role
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, MapDBError(err, "GetAllUsers")
	}
	defer rows.Close()

	var users []*pb.User

	for rows.Next() {
		var id, name string
		var telegramID int64
		var role string

		err := rows.Scan(&id, &name, &telegramID, &role)

		if err != nil {
			return nil, MapDBError(err, "GetAllUsers")
		}

		user := &pb.User{
			Id:         id,
			Name:       name,
			TelegramId: telegramID,
			Role:       role,
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, MapDBError(err, "GetAllUsers")
	}

	return users, nil
}

func (r *Repository) UpdateTutorBio(userID string, bio string) error {
	const query = `
		UPDATE tutors SET
            bio = $1 
        WHERE id = $2`

	_, err := r.db.Exec(query, bio, userID)

	if err != nil {
		return MapDBError(err, "UpdateTutorBio")
	}

	return nil
}

func (r *Repository) UpdateTutorTags(tutorID string, tags []string) error {
	const query = `
		UPDATE tutors SET 
		    tags = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, pq.Array(tags), tutorID)

	if err != nil {
		return MapDBError(err, "UpdateTutorTags")
	}

	return nil
}

func (r *Repository) SetNewIsActiveTutor(tutorID string, isActive bool) error {
	const query = `
		UPDATE tutors SET
            is_active = $1
        WHERE id = $2`

	_, err := r.db.Exec(query, isActive, tutorID)

	if err != nil {
		return MapDBError(err, "UpdateTutorTags")
	}

	return nil
}

func (r *Repository) UpdateTutorName(tutorID string, name string) error {
	const query = `
		UPDATE users SET
		    name = $1
		WHERE id = $2 AND role = $3`

	_, err := r.db.Exec(query, name, tutorID, models.RoleTutor)

	if err != nil {
		return MapDBError(err, "UpdateTutorTags")
	}

	return nil
}

func (r *Repository) GetAllTutorsPagination(limit int, offset int, tag string) ([]*pb.Tutor, int, error) {
	queryCount := `
		SELECT 
			COUNT(*)
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true AND u.is_banned = false`

	argsCount := []interface{}{models.RoleTutor}

	if tag != "" {
		tag = strings.ToLower(tag)
		queryCount += ` AND $2 = ANY(t.tags)`
		argsCount = append(argsCount, tag)
	}

	var total int

	err := r.db.QueryRow(queryCount, argsCount...).Scan(&total)

	if err != nil {
		custError := MapDBError(err, "GetAllTutorsPagination")
		return nil, 0, custError
	}

	queryGetAllPagination := `
		SELECT 
			u.id, 
			u.name, 
			u.telegram_id,
			t.rating,
			t.tags
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true AND u.is_banned = false`

	argsPagination := []interface{}{models.RoleTutor, limit, offset}

	if tag != "" {
		tag = strings.ToLower(tag)
		queryGetAllPagination += ` AND $4 = ANY(t.tags)`
		argsPagination = append(argsPagination, tag)
	}

	queryGetAllPagination += ` ORDER BY t.rating DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(queryGetAllPagination, argsPagination...)

	if err != nil {
		custError := MapDBError(err, "GetAllTutorsPagination")
		return nil, 0, custError
	}
	defer rows.Close()

	var tutors []*pb.Tutor

	for rows.Next() {
		var id string
		var name string
		var telegramID int64
		var rating int32
		var tags pq.StringArray

		err := rows.Scan(&id, &name, &telegramID, &rating, &tags)
		if err != nil {
			custom_logger.Info(fmt.Sprintf("[Postgres] GetAllTutorsPagination failed. Error: %v", err.Error()))
			return nil, 0, err
		}

		tutor := &pb.Tutor{
			User: &pb.User{
				Id:         id,
				Name:       name,
				TelegramId: telegramID,
			},
			Tags:   tags,
			Rating: rating,
		}
		tutors = append(tutors, tutor)
	}

	err = rows.Err()

	if err != nil {
		custError := MapDBError(err, "GetAllTutorsPagination")
		return nil, 0, custError
	}

	return tutors, total, nil
}

func (r *Repository) CreateReview(tutorID, orderID string, rating int, comment string) (string, error) {
	timestamp := time.Now()

	queryInsertReview := `
		INSERT INTO reviews (
            tutor_id,
            order_id,
            rating,
            comment,
		    is_active,
            created_at
        )
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var reviewID string
	err := r.db.QueryRow(queryInsertReview, tutorID, orderID, rating, comment, false, timestamp).Scan(&reviewID)

	if err != nil {
		return "", MapDBError(err, "CreateReview")
	}

	return reviewID, nil
}

func (r *Repository) GetReviews(tutorID string) ([]models.Review, error) {
	const query = `
		SELECT 
			id, 
			tutor_id, 
			order_id, 
			rating, 
			comment, 
			is_active,
			created_at
		FROM reviews
		WHERE tutor_id = $1`

	var reviews []models.Review
	err := r.db.Select(&reviews, query, tutorID)

	if err != nil {
		return nil, MapDBError(err, "CreateReview")
	}

	return reviews, nil
}

func (r *Repository) GetReviewById(reviewID string) (*models.Review, error) {
	const query = `
		SELECT 
			id,
			tutor_id,
			order_id,
			rating,
			comment,
			is_active,
			created_at
		FROM reviews
		WHERE id = $1`

	var review models.Review
	err := r.db.Get(&review, query, reviewID)

	if err != nil {
		return nil, MapDBError(err, "GetReviewById")
	}

	return &review, nil
}

func (r *Repository) RemoveOneResponse(tutorID string) error {
	const query = `
		UPDATE tutors SET
            response_count = response_count - 1
        WHERE id = $1 AND response_count > 0
        RETURNING response_count`

	var newCount int64
	err := r.db.QueryRow(query, tutorID).Scan(&newCount)

	if err != nil {
		return MapDBError(err, "RemoveOneResponse")
	}

	return nil
}

func (r *Repository) AddResponses(tutorTelegramID int64, responseCount int) (int, error) {
	const query = `
        UPDATE tutors
		SET response_count = response_count + $1
		WHERE id = (
			SELECT id FROM users WHERE telegram_id = $2 AND role = $3
		)
		RETURNING response_count`

	var newCount int64
	err := r.db.QueryRow(query, responseCount, tutorTelegramID, models.RoleTutor).Scan(&newCount)

	if err != nil {
		return 0, MapDBError(err, "AddResponses")
	}

	return int(newCount), nil
}

func (r *Repository) SetReviewActive(reviewID, tutorID string) error {
	tx, err := r.db.Begin()

	if err != nil {
		return MapDBError(err, "SetReviewActive")
	}
	defer tx.Rollback()

	const querySetActive = `
		UPDATE reviews SET
		    is_active = true
		WHERE id = $1`

	res, err := tx.Exec(querySetActive, reviewID)

	if err != nil {
		return MapDBError(err, "SetReviewActive")
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return MapDBError(err, "SetReviewActive")
	}

	if rowsAffected == 0 {
		custom_logger.Info("[Postgres] SetReviewActive no review found with ID: " + reviewID)
		return custom_errors.ErrorNotFound
	}

	const queryGetAvgRating = `
		SELECT
		    ROUND(AVG(rating))::INT
		FROM reviews
		WHERE tutor_id = $1 AND is_active = true`

	var rating int32

	err = tx.QueryRow(queryGetAvgRating, tutorID).Scan(&rating)

	if err != nil {
		return MapDBError(err, "SetReviewActive")
	}

	const queryUpdateTutor = `
		UPDATE tutors SET
		    rating = $1
		WHERE id = $2`

	_, err = tx.Exec(queryUpdateTutor, rating, tutorID)

	if err != nil {
		return MapDBError(err, "SetReviewActive")
	}

	err = tx.Commit()

	if err != nil {
		return MapDBError(err, "SetReviewActive")
	}

	return nil
}

func (r *Repository) GetAllTutorsResponseCondition(minResponseCount int) ([]*models.TutorWithResponse, error) {
	custom_logger.Info("[Postgres] GetAllTutorsResponseCondition called")

	const query = `
		SELECT 
			u.id, 
			u.name, 
			u.telegram_id, 
			t.response_count
		FROM users u 
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.response_count < $2`

	rows, err := r.db.Query(query, models.RoleTutor, minResponseCount)

	if err != nil {
		return nil, MapDBError(err, "GetAllTutorsResponseCondition")
	}
	defer rows.Close()

	var tutors []*models.TutorWithResponse

	for rows.Next() {
		var id, name string
		var telegramID int64
		var responseCount int32

		err := rows.Scan(&id, &name, &telegramID, &responseCount)

		if err != nil {
			return nil, MapDBError(err, "GetAllTutorsResponseCondition")
		}

		tutor := &models.TutorWithResponse{
			Id:            id,
			Name:          name,
			TelegramID:    telegramID,
			ResponseCount: responseCount,
		}
		tutors = append(tutors, tutor)
	}

	err = rows.Err()

	if err != nil {
		return nil, MapDBError(err, "GetAllTutorsResponseCondition")
	}

	return tutors, nil
}

func (r *Repository) BanUser(telegramID int64, isBanned bool) error {
	const query = `
		UPDATE users SET
            is_banned = $1
        WHERE telegram_id = $2`

	_, err := r.db.Exec(query, isBanned, telegramID)

	if err != nil {
		return MapDBError(err, "GetAllTutorsResponseCondition")
	}

	return nil
}
