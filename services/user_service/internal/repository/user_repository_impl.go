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
	pb "github.com/randnull/Lessons/internal/gRPC"
	lg "github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	"log"
	"strconv"
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
		log.Fatal("[Postgres] failed to connect" + err.Error())
	}

	err = db.PingContext(context.Background())

	if err != nil {
		log.Fatal("[Postgres] failed to ping" + err.Error())
	}

	log.Print("[Postgres] Database is ready")

	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(user *models.CreateUser) (string, error) {
	lg.Info(fmt.Sprintf("[Postgres] CreateUser called. UserID: %v. Role: %v", user.TelegramId, user.Role))

	ExistedUser, err := r.GetUserByTelegramId(user.TelegramId, user.Role)

	if err == nil {
		lg.Info(fmt.Sprintf("[Postgres] CreateUser User exist. Authorization Telegram-UserID: %v. UserID: %v", user.TelegramId, ExistedUser.Id))
		return ExistedUser.Id, nil
	}

	tx, err := r.db.Begin()
	defer tx.Rollback()

	if err != nil {
		lg.Info("[Postgres] CreateUser TX failed " + fmt.Sprint(user.TelegramId) + " Error: " + err.Error())
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
		lg.Info("[Postgres] CreateUser failed: " + err.Error())
		return "", custom_errors.ErrorWithCreate
	}

	if user.Role == models.RoleTutor {
		const queryInsertTutor = `
				INSERT INTO tutors (
                    id,
                    created_at)
				VALUES ($1, $2)`

		_, err = tx.Exec(queryInsertTutor, UserId, currentTime)

		if err != nil {
			lg.Info("[Postgres] CreateUser failed: " + err.Error())
			tx.Rollback()
			return "", custom_errors.ErrorWithCreate
		}
	}

	err = tx.Commit()
	if err != nil {
		lg.Info("[Postgres] CreateUser failed: " + err.Error())
		return "", err
	}

	return UserId, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		lg.Error("[Postgres] GetUserByTelegramId error: " + err.Error())
		return nil, err
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		lg.Error("[Postgres] GetStudentById error: " + err.Error())
		return nil, err
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		lg.Error("[Postgres] GetTutorByID failed. UserID: " + userID + " Error: " + err.Error())
		return nil, err
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
		lg.Error("[Postgres] GetAllUsers failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User

	for rows.Next() {
		var id, name string
		var telegramID int64
		var role string

		err := rows.Scan(&id, &name, &telegramID, &role)

		if err != nil {
			lg.Error("[Postgres] GetAllTutors failed: " + err.Error())
			return nil, err
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
		lg.Error("[Postgres] GetAllTutors error after rows scan")
		return nil, custom_errors.ErrorAfterRowScan
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
		lg.Error("[Postgres] UpdateTutorBio failed. Error: " + err.Error())
		return err
	}

	return nil
}

func (r *Repository) UpdateTutorTags(tutorID string, tags []string) error {
	lg.Info("[Postgres] UpdateTutorTags called. TutorID: " + tutorID)

	const query = `
		UPDATE tutors SET 
		    tags = $1
		WHERE id = $2`

	_, err := r.db.Exec(query, pq.Array(tags), tutorID)

	if err != nil {
		lg.Error("[Postgres] UpdateTutorTags failed. Error: " + err.Error())
		return err
	}

	lg.Info("[Postgres] UpdateTutorTags success")
	return nil
}

func (r *Repository) SetNewIsActiveTutor(tutorID string, isActive bool) error {
	lg.Info("[Postgres] SetNewIsActiveTutor called. TutorID: " + tutorID)

	const query = `
		UPDATE tutors SET
            is_active = $1
        WHERE id = $2`

	_, err := r.db.Exec(query, isActive, tutorID)

	if err != nil {
		lg.Error("[Postgres] SetNewIsActiveTutor failed. Error: " + err.Error())
		return err
	}

	lg.Info("[Postgres] SetNewIsActiveTutor success")
	return nil
}

func (r *Repository) UpdateTutorName(tutorID string, name string) error {
	lg.Info("[Postgres] UpdateTutorName called. TutorID: " + tutorID + " Name: " + name)

	const query = `
		UPDATE users SET
		    name = $1
		WHERE id = $2 AND role = $3`

	_, err := r.db.Exec(query, name, tutorID, models.RoleTutor)

	if err != nil {
		lg.Error("[Postgres] UpdateTutorName failed. Error: " + err.Error())
		return err
	}

	lg.Info("[Postgres] UpdateTutorName success")
	return nil
}

func (r *Repository) GetAllTutorsPagination(limit int, offset int, tag string) ([]*pb.Tutor, int, error) {
	const queryCount = `
		SELECT 
			COUNT(*)
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true AND u.is_banned = false`

	var total int
	err := r.db.QueryRow(queryCount, models.RoleTutor).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	lg.Info("[Postgres] GetAllTutorsPagination called.")

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
		lg.Info(fmt.Sprintf("[Postgres] GetAllTutorsPagination failed. Error: %v", err.Error()))
		return nil, 0, err
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
			lg.Info(fmt.Sprintf("[Postgres] GetAllTutorsPagination failed. Error: %v", err.Error()))
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
		lg.Info(fmt.Sprintf("[Postgres] GetAllTutorsPagination failed. Error: %v", err.Error()))
		return nil, 0, err
	}

	lg.Info("[Postgres] GetAllTutorsPagination success.")

	return tutors, total, nil
}

func (r *Repository) CreateReview(tutorID, orderID string, rating int, comment string) (string, error) {
	lg.Info("[Postgres] CreateReview called. TutorID: " + tutorID + " orderID: " + orderID)

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
		lg.Info("[Postgres] CreateReview failed. TutorID: " + tutorID + " orderID: " + orderID + " Error: " + err.Error())
		return "", err
	}

	lg.Info("[Postgres] CreateReview success. TutorID: " + tutorID + " orderID: " + orderID)

	return reviewID, nil
}

func (r *Repository) GetReviews(tutorID string) ([]models.Review, error) {
	lg.Info("[Postgres] GetReviews called. tutorID: " + tutorID)

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
		if errors.Is(err, sql.ErrNoRows) {
			lg.Info(fmt.Sprintf("[Postgres] GetReviews. tutorID: %v. Not found", tutorID))
			return nil, custom_errors.UserNotFound
		}
		lg.Error(fmt.Sprintf("[Postgres] GetReviewById failed. tutorID: %v. Error: %v", tutorID, err.Error()))
		return nil, custom_errors.ErrorServiceError
	}

	lg.Info("[Postgres] GetReviews success. tutorID: " + tutorID)

	return reviews, nil
}

func (r *Repository) GetReviewById(reviewID string) (*models.Review, error) {
	lg.Info("[Postgres] GetReviewById called. reviewID: " + reviewID)

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
		if errors.Is(err, sql.ErrNoRows) {
			lg.Info("[Postgres] GetReviewById. ReviewID: " + fmt.Sprint(reviewID) + " Not found.")
			return nil, custom_errors.UserNotFound
		}
		lg.Error("[Postgres] GetReviewById failed. ReviewID: " + fmt.Sprint(reviewID) + " Not found.")
		return nil, custom_errors.ErrorServiceError
	}

	lg.Info("[Postgres] GetReviewById success. reviewID: " + reviewID)

	return &review, nil
}

func (r *Repository) GetTagsByTutorID(tutorID string) ([]string, error) {
	lg.Info("[Postgres] GetTagsByTutorID called. tutorID: " + tutorID)

	const query = `
		SELECT
    		tags
		FROM tutors WHERE id = $1`

	rows, err := r.db.Query(query, tutorID)

	if err != nil {
		lg.Error("[Postgres] GetTagsByTutorID failed.  tutorID: " + tutorID + " Error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagArray []string
		err := rows.Scan(pq.Array(&tagArray))
		if err != nil {
			lg.Error("[Postgres] GetTagsByTutorID failed. tutorID: " + tutorID + " Error: " + err.Error())
			return nil, err
		}
		tags = append(tags, tagArray...)
	}

	lg.Info("[Postgres] GetTagsByTutorID success. tutorID: " + tutorID)
	return tags, nil
}

func (r *Repository) RemoveOneResponse(tutorID string) error {
	lg.Info("[Postgres] RemoveOneResponse called. tutorID: " + tutorID)

	const query = `
		UPDATE tutors SET
            response_count = response_count - 1
        WHERE id = $1 AND response_count > 0
        RETURNING response_count`

	var newCount int64
	err := r.db.QueryRow(query, tutorID).Scan(&newCount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return custom_errors.ErrorNotFound
		}
		lg.Info("[Postgres] RemoveOneResponse failed. tutorID: " + tutorID + " Error: " + err.Error())
		return err
	}

	lg.Info("[Postgres] RemoveOneResponse success. tutorID: " + tutorID)
	return nil
}

func (r *Repository) AddResponses(tutorTelegramID int64, responseCount int) (int, error) {
	lg.Info("[Postgres] AddResponses called. tutorTelegramID: " + fmt.Sprint(tutorTelegramID) + " ResponseCount: " + strconv.Itoa(responseCount))

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
		if errors.Is(err, sql.ErrNoRows) {
			lg.Error("[Postgres] AddResponses failed. tutorTelegramID: " + fmt.Sprint(tutorTelegramID) + " Not found.")
			return 0, custom_errors.UserNotFound
		}
		lg.Error("[Postgres] AddResponses failed. tutorTelegramID: " + fmt.Sprint(tutorTelegramID) + " Error: " + err.Error())
		return 0, err
	}

	lg.Info("[Postgres] AddResponses success. tutorTelegramID: " + fmt.Sprint(tutorTelegramID) + " NewCount: " + strconv.Itoa(int(newCount)))
	return int(newCount), nil
}

func (r *Repository) SetReviewActive(reviewID, tutorID string) error {
	lg.Info("[Postgres] SetReviewActive called. ReviewID: " + reviewID)

	tx, err := r.db.Begin()

	if err != nil {
		lg.Error("[Postgres] SetReviewActive failed: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const querySetActive = `
		UPDATE reviews SET
		    is_active = true
		WHERE id = $1`

	res, err := tx.Exec(querySetActive, reviewID)
	if err != nil {
		lg.Error("[Postgres] SetReviewActive failed: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		lg.Error("[Postgres] SetReviewActive failed: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	if rowsAffected == 0 {
		lg.Info("[Postgres] SetReviewActive no review found with ID: " + reviewID)
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
		lg.Error("[Postgres] SetReviewActive failed: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	const queryUpdateTutor = `
		UPDATE tutors SET
		    rating = $1
		WHERE id = $2`

	_, err = tx.Exec(queryUpdateTutor, rating, tutorID)

	if err != nil {
		lg.Error("[Postgres] SetReviewActive failed queryUpdateTutor: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	if err = tx.Commit(); err != nil {
		lg.Error("[Postgres] SetReviewActive failed: " + err.Error())
		return custom_errors.ErrorServiceError
	}

	lg.Info("[Postgres] SetReviewActive success")
	return nil
}

func (r *Repository) GetAllTutorsResponseCondition(minResponseCount int) ([]*models.TutorWithResponse, error) {
	lg.Info("[Postgres] GetAllTutorsResponseCondition called")

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
		lg.Error("[Postgres] GetAllTutorsResponseCondition failed. Query error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tutors []*models.TutorWithResponse

	for rows.Next() {
		var id, name string
		var telegramID int64
		var responseCount int32

		err := rows.Scan(&id, &name, &telegramID, &responseCount)

		if err != nil {
			lg.Error("[Postgres] GetAllTutorsResponseCondition error: " + err.Error())
			return nil, err
		}

		tutor := &models.TutorWithResponse{
			Id:            id,
			Name:          name,
			TelegramID:    telegramID,
			ResponseCount: responseCount,
		}
		tutors = append(tutors, tutor)
	}

	if rows.Err() != nil {
		lg.Error("[Postgres] GetAllTutorsResponseCondition error after rows scan")
		return nil, custom_errors.ErrorAfterRowScan
	}

	lg.Info("[Postgres] GetAllTutorsResponseCondition success")

	return tutors, nil
}

func (r *Repository) BanUser(telegramID int64, isBanned bool) error {
	const query = `
		UPDATE users SET
            is_banned = $1
        WHERE telegram_id = $2`

	_, err := r.db.Exec(query, isBanned, telegramID)

	if err != nil {
		lg.Error("[Postgres] BanUser failed. Error: " + err.Error())
		return err
	}

	return nil
}

func (r *Repository) GetUserById(userID string) (*models.UserDB, error) {
	lg.Info("[Postgres] GetUserById called. UserID: " + userID)

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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		lg.Error("[Postgres] GetUserById error: " + err.Error())
		return nil, err
	}

	lg.Info("[Postgres] GetUserById success. UserID: " + userID)

	return &user, nil
}
