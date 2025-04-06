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
	"github.com/randnull/Lessons/internal/models"
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

func (r *Repository) CreateUser(user *models.CreateUser) (string, error) {
	if (user.Role != models.RoleStudent) && (user.Role != models.RoleTutor) {
		return "", custom_errors.ErrorIncorrectRole
	}

	ExistedUser, err := r.GetUserByTelegramId(user.TelegramId, user.Role)
	if err == nil {
		log.Println("exist")
		return ExistedUser.Id, nil
	}

	tx, err := r.db.Begin()
	defer tx.Rollback()

	if err != nil {
		log.Println("err:", err)
		tx.Rollback()
		return "", err
	}

	query := `INSERT INTO users (telegram_id, name, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id`

	var UserId string

	currentTime := time.Now()

	err = tx.QueryRow(query,
		user.TelegramId,
		user.Name,
		user.Role,
		currentTime,
	).Scan(&UserId)

	log.Println(err)
	if err != nil {
		log.Println("err 2:", err)
		tx.Rollback()
		return "", custom_errors.ErrorWithCreate
	}

	if user.Role == models.RoleTutor {
		queryInsertTutor := `INSERT INTO tutors (id, created_at) VALUES ($1, $2)`
		_, err = tx.Exec(queryInsertTutor, UserId, currentTime)

		if err != nil {
			log.Println("err 3:", err)

			tx.Rollback()
			return "", custom_errors.ErrorWithCreate
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)

		return "", err
	}

	return UserId, nil
}

func (r *Repository) GetUserByTelegramId(telegramID int64, userRole string) (*models.UserDB, error) {
	user := &models.UserDB{}

	query := `SELECT id, telegram_id, name, role, created_at FROM users WHERE telegram_id = $1 AND role = $2`

	err := r.db.QueryRow(query, telegramID, userRole).Scan(
		&user.Id,
		&user.TelegramID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		log.Println('0', err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetStudentById(userID string) (*models.UserDB, error) {
	user := &models.UserDB{}

	query := `SELECT id, telegram_id, name, role, created_at FROM users WHERE id = $1 AND role = $2`

	err := r.db.QueryRow(query, userID, "Student").Scan(
		&user.Id,
		&user.TelegramID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetTutorByID(userID string) (*models.TutorDB, error) {
	log.Println("come here 3")

	user := &models.TutorDB{}

	query := `SELECT 
    			u.id, 
    			u.telegram_id, 
    			u.name, 
    			u.role, 
    			u.created_at,
    			t.bio
			FROM users u 
			LEFT JOIN tutors t ON u.id = t.id
			WHERE u.id = $1 AND u.role = $2`

	log.Println("come here 5")

	var bio sql.NullString

	err := r.db.QueryRow(query, userID, "Tutor").Scan(
		&user.Id,
		&user.TelegramID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&bio,
	)
	if bio.Valid {
		user.Bio = bio.String
	}

	log.Println(err)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserById(userID string) (*models.UserDB, error) {
	user := &models.UserDB{}

	query := `SELECT id, telegram_id, name, role, created_at FROM users WHERE id = $1 AND role = $2`

	err := r.db.QueryRow(query, userID, "Tutor").Scan(
		&user.Id,
		&user.TelegramID,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, custom_errors.UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetAllTutors() ([]*pb.User, error) {
	query := `SELECT id, name FROM users WHERE role = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, "Tutor")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var Users []*pb.User

	for rows.Next() {
		var user pb.User

		err := rows.Scan(
			&user.Id,
			&user.Name,
		)

		if err != nil {
			return nil, err
		}

		Users = append(Users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Users, nil
}

func (r *Repository) UpdateTutorBio(userID string, bio string) error {
	queryUpdateBioTutor := `UPDATE tutors SET bio = $1 WHERE id = $2`

	_, err := r.db.Exec(queryUpdateBioTutor, bio, userID)

	if err != nil {
		log.Println(err)
		return custom_errors.ErrorUpdateBio
	}

	return nil
}

func (r *Repository) GetAllTutorsPagination(limit int, offset int) ([]*pb.User, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	var total int

	queryCount := `SELECT 
    					COUNT(*) 
					FROM users WHERE role = $1`

	err = tx.QueryRow(queryCount, "Tutor").Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	queryGetAllPagination := `SELECT
    							id, 
    							name 
							FROM users WHERE role = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := tx.Query(queryGetAllPagination, "Tutor", limit, offset)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}

	var Users []*pb.User

	for rows.Next() {
		var user pb.User

		err := rows.Scan(
			&user.Id,
			&user.Name,
		)

		if err != nil {
			return nil, 0, err
		}

		Users = append(Users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	if err = tx.Commit(); err != nil {
		return nil, 0, err
	}

	return Users, total, nil
}

func (r *Repository) UpdateTutorTags(tutorID string, tags []string) error {
	queryUpdateTutorTags := `
		UPDATE tutors
		SET tags = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(queryUpdateTutorTags, pq.Array(tags), tutorID)
	if err != nil {
		log.Printf("Failed to update tags for tutor")
		return err //custom_errors.ErrorTagsTutor
	}

	return nil
}

func (r *Repository) CreateReview(tutorID, studentID string, rating int, comment string) (string, error) {
	timestamp := time.Now()

	queryInsertReview := `INSERT INTO reviews (
                     tutor_id,
                     student_id,
                     rating,
                     comment,
                     created_at
        )
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var reviewID string
	err := r.db.QueryRow(queryInsertReview, tutorID, studentID, rating, comment, timestamp).Scan(&reviewID)
	if err != nil {
		log.Printf("Failed create review ((: %v", err)
		return "", err
	}

	return reviewID, nil
}

func (r *Repository) GetReviews(tutorID string) ([]models.Review, error) {
	query := `SELECT 
    			id, 
    			tutor_id, 
    			student_id, 
    			rating, 
    			comment, 
    			created_at
		FROM reviews
		WHERE tutor_id = $1
	`

	rows, err := r.db.Query(query, tutorID)
	if err != nil {
		log.Printf("Failed get reviews for tutor %s: %v", tutorID, err)
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review

	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.TutorID, &review.StudentID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			log.Printf("Failed to scan review row: %v", err)
			continue
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	return reviews, nil
}

func (r *Repository) GetReviewById(reviewID string) (*models.Review, error) {
	query := `SELECT 
    			id,
       			tutor_id,
       			student_id,
       			rating,
       			comment,
       			created_at
		FROM reviews
		WHERE id = $1
	`

	var review models.Review
	err := r.db.QueryRow(query, reviewID).Scan(
		&review.ID, &review.TutorID, &review.StudentID, &review.Rating, &review.Comment, &review.CreatedAt,
	)
	if err != nil {
		log.Printf("Failed get review with ID %s: %v", reviewID, err)
		return nil, err
	}

	return &review, nil
}

func (r *Repository) GetTagsByTutorID(tutorID string) ([]string, error) {
	query := `SELECT
    			tags
			FROM tutors WHERE id = $1`
	rows, err := r.db.Query(query, tutorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagArray []string
		if err := rows.Scan(pq.Array(&tagArray)); err != nil {
			return nil, err
		}
		tags = append(tags, tagArray...)
	}
	return tags, nil
}

func (r *Repository) SetNewIsActiveTutor(tutorID string, isActive bool) error {
	query := `UPDATE tutors SET
                is_active = $1
              WHERE id = $2`

	_, err := r.db.Exec(query, isActive, tutorID)
	if err != nil {
		log.Printf("Failed is_active for tutor %s: %v", tutorID, err)
		return err
	}

	return nil
}
