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
	lg.Info("[Postgres] CreateUser called. UserID: " + fmt.Sprint(user.TelegramId) + " Role: " + user.Role)

	ExistedUser, err := r.GetUserByTelegramId(user.TelegramId, user.Role)

	if err == nil {
		lg.Info("[Postgres] CreateUser User exist. Authorization UserID: " + fmt.Sprint(user.TelegramId) + " UserID " + ExistedUser.Id)
		return ExistedUser.Id, nil
	}

	tx, err := r.db.Begin()
	defer tx.Rollback()

	if err != nil {
		lg.Info("[Postgres] CreateUser TX failed " + fmt.Sprint(user.TelegramId) + " Error: " + err.Error())
		tx.Rollback()
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
		lg.Info("[Postgres] CreateUser Insert failed for userTelegramID: " + fmt.Sprint(user.TelegramId) + " Error: " + err.Error())
		tx.Rollback()
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
			lg.Info("[Postgres] CreateUser Insert Tutor failed for userTelegramID: " + fmt.Sprint(user.TelegramId) + " Error: " + err.Error())

			tx.Rollback()
			return "", custom_errors.ErrorWithCreate
		}
	}

	err = tx.Commit()
	if err != nil {
		lg.Info("[Postgres] CreateUser TX failed " + fmt.Sprint(user.TelegramId) + " Error: " + err.Error())
		return "", err
	}

	return UserId, nil
}

func (r *Repository) GetUserByTelegramId(telegramID int64, userRole string) (*models.UserDB, error) {
	lg.Info("[Postgres] GetUserByTelegramId called. UserID: " + fmt.Sprint(telegramID) + " Role: " + userRole)

	var user models.UserDB

	const query = `
		SELECT 
		    id,
		    telegram_id,
		    name, 
		    role, 
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

	lg.Info("[Postgres] GetUserByTelegramId success. UserTelegramID: " + fmt.Sprint(telegramID) + " Role: " + userRole + " UserID: " + user.Id)

	return &user, nil
}

func (r *Repository) GetStudentById(userID string) (*models.UserDB, error) {
	lg.Info("[Postgres] GetStudentById called. UserID: " + userID)

	var user models.UserDB

	const query = `
		SELECT 
		    id,
		    telegram_id,
		    name, 
		    role, 
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

	lg.Info("[Postgres] GetStudentById success. UserID: " + userID)

	return &user, nil
}

func (r *Repository) GetTutorByID(userID string) (*models.TutorDB, error) {
	lg.Info("[Postgres] GetTutorByID called. UserID: " + userID)

	const query = `
        SELECT 
            u.id, 
            u.telegram_id, 
            u.name, 
            u.role, 
            u.created_at,
            t.bio,
            t.response_count,
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

	err := r.db.QueryRow(query, userID, models.RoleTutor).Scan(
		&tutor.Id,
		&tutor.TelegramID,
		&tutor.Name,
		&tutor.Role,
		&tutor.CreatedAt,
		&bio,
		&responseCount,
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

	lg.Info("[Postgres] GetTutorByID succss. UserID: " + userID)

	return &tutor, nil
}

//func (r *Repository) GetUserById(userID string) (*models.UserDB, error) {
//	user := &models.UserDB{}
//
//	query := `SELECT id, telegram_id, name, role, created_at FROM users WHERE id = $1 AND role = $2`
//
//	err := r.db.QueryRow(query, userID, "Tutor").Scan(
//		&user.Id,
//		&user.TelegramID,
//		&user.Name,
//		&user.Role,
//		&user.CreatedAt,
//	)
//
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, custom_errors.UserNotFound
//		}
//		return nil, err
//	}
//
//	return user, nil
//}

func (r *Repository) GetAllTutors() ([]*pb.Tutor, error) {
	lg.Info("[Postgres] GetAllTutors called")

	const query = `
		SELECT 
			u.id, 
			u.name, 
			u.telegram_id, 
			t.tags
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true
		ORDER BY u.created_at DESC`

	rows, err := r.db.Query(query, models.RoleTutor)
	if err != nil {
		lg.Error("[Postgres] GetAllTutors failed. Query error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tutors []*pb.Tutor

	for rows.Next() {
		var id, name string
		var telegramID int64
		var tags pq.StringArray

		err := rows.Scan(&id, &name, &telegramID, &tags)

		if err != nil {
			lg.Error("[Postgres] GetAllTutors error: " + err.Error())
			return nil, err
		}

		tutor := &pb.Tutor{
			User: &pb.User{
				Id:         id,
				Name:       name,
				TelegramId: telegramID,
			},
			Tags: tags,
		}
		tutors = append(tutors, tutor)
	}

	if rows.Err() != nil {
		lg.Error("[Postgres] GetAllTutors error after rows scan")
		return nil, custom_errors.ErrorAfterRowScan
	}

	lg.Info("[Postgres] GetAllTutors success")

	return tutors, nil
}

func (r *Repository) UpdateTutorBio(userID string, bio string) error {
	lg.Info("[Postgres] UpdateTutorBio called")

	const query = `
		UPDATE tutors SET
            bio = $1 
        WHERE id = $2`

	_, err := r.db.Exec(query, bio, userID)

	if err != nil {
		lg.Error("[Postgres] UpdateTutorBio failed. Error: " + err.Error())
		return err
	}

	lg.Info("[Postgres] UpdateTutorBio success")
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

// этого монстра нужно отрефакторить
func (r *Repository) GetAllTutorsPagination(limit int, offset int, tag string) ([]*pb.Tutor, int, error) {
	var total int

	queryCount := `
		SELECT
		    COUNT(*) 
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true`

	argsCount := []interface{}{"Tutor"}

	if tag != "" {
		queryCount += ` AND $2 = ANY(t.tags)`
		argsCount = append(argsCount, tag)
	}

	err := r.db.QueryRow(queryCount, argsCount...).Scan(&total)

	if err != nil {
		return nil, 0, errors.New("error with rows")
	}

	queryGetAllPagination := `
		SELECT 
			u.id, 
			u.name, 
			u.telegram_id,
			t.tags
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true`

	argsPagination := []interface{}{"Tutor", limit, offset}

	if tag != "" {
		tag = strings.ToLower(tag)
		queryGetAllPagination += ` AND $4 = ANY(t.tags)`
		argsPagination = append(argsPagination, tag)
	}

	queryGetAllPagination += ` ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(queryGetAllPagination, argsPagination...)

	if err != nil {
		return nil, 0, errors.New("error with rows")
	}
	defer rows.Close()

	var tutors []*pb.Tutor

	for rows.Next() {
		var id string
		var name string
		var telegramID int64
		var tags pq.StringArray

		err := rows.Scan(&id, &name, &telegramID, &tags)
		if err != nil {
			return nil, 0, errors.New("error with scan")
		}

		tutor := &pb.Tutor{
			User: &pb.User{
				Id:         id,
				Name:       name,
				TelegramId: telegramID,
			},
			Tags: tags,
		}
		tutors = append(tutors, tutor)
	}

	err = rows.Err()

	if err != nil {
		return nil, 0, errors.New("error with rows")
	}

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
		WHERE tutor_id = $1 AND is_active = $2`

	var reviews []models.Review
	err := r.db.Select(&reviews, query, tutorID, true)

	if err != nil {
		lg.Error("[Postgres] GetReviews failed. tutorID:" + tutorID + " Error: " + err.Error())
		return nil, err
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
		lg.Error("[Postgres] GetReviewById failed. reviewID:" + reviewID + " Error: " + err.Error())
		return nil, err
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

func (r *Repository) SetReviewActive(reviewID string) error {
	lg.Info("[Postgres] SetReviewActive called. ReviewID: " + reviewID)

	const query = `
		UPDATE reviews SET
            is_active = $1
        WHERE id = $2`

	_, err := r.db.Exec(query, true, reviewID)

	if err != nil {
		lg.Error("[Postgres] SetReviewActive failed. Error: " + err.Error())
		return err
	}

	lg.Info("[Postgres] SetReviewActive success")
	return nil
}
