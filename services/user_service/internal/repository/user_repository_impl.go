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
	query := `
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

	err := r.db.QueryRow(query, userID, "Tutor").Scan(
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
		return nil, fmt.Errorf("failed to query tutor: %w", err)
	}

	if bio.Valid {
		tutor.Bio = bio.String
	}

	if responseCount.Valid {
		tutor.ResponseCount = int32(responseCount.Int32)
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

	return &tutor, nil
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

func (r *Repository) GetAllTutors() ([]*pb.Tutor, error) {
	query := `
		SELECT 
			u.id, 
			u.name, 
			u.telegram_id, 
			t.tags
		FROM users u
		JOIN tutors t ON u.id = t.id
		WHERE u.role = $1 AND t.is_active = true
		ORDER BY u.created_at DESC`

	rows, err := r.db.Query(query, "Tutor")
	if err != nil {
		return nil, errors.New("error with query")
	}
	defer rows.Close()

	var tutors []*pb.Tutor
	for rows.Next() {
		var id, name string
		var telegramID int64
		var tags pq.StringArray

		err := rows.Scan(&id, &name, &telegramID, &tags)
		if err != nil {
			return nil, errors.New("error with scan")
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

	if err = rows.Err(); err != nil {
		return nil, errors.New("error with rows")
	}

	return tutors, nil
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

func (r *Repository) GetAllTutorsPagination(limit int, offset int, tag string) ([]*pb.Tutor, int, error) {
	var total int

	queryCount := `
		SELECT COUNT(*) 
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

func (r *Repository) CreateReview(tutorID, orderID string, rating int, comment string) (string, error) {
	timestamp := time.Now()

	queryInsertReview := `INSERT INTO reviews (
                     tutor_id,
                     order_id,
                     rating,
                     comment,
                     created_at
        )
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var reviewID string
	err := r.db.QueryRow(queryInsertReview, tutorID, orderID, rating, comment, timestamp).Scan(&reviewID)
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
    			order_id, 
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
		err := rows.Scan(&review.ID, &review.TutorID, &review.OrderID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		reviews = append(reviews, review)
	}

	err = rows.Err()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return reviews, nil
}

func (r *Repository) GetReviewById(reviewID string) (*models.Review, error) {
	query := `SELECT 
    			id,
       			tutor_id,
       			order_id,
       			rating,
       			comment,
       			created_at
		FROM reviews
		WHERE id = $1
	`

	var review models.Review
	err := r.db.QueryRow(query, reviewID).Scan(
		&review.ID, &review.TutorID, &review.OrderID, &review.Rating, &review.Comment, &review.CreatedAt,
	)

	if err != nil {
		log.Println(err)
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
		err := rows.Scan(pq.Array(&tagArray))
		if err != nil {
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
		log.Println(err)
		return err
	}

	return nil
}

func (r *Repository) RemoveOneResponse(tutorID string) error {
	query := `UPDATE tutors SET
                response_count = response_count - 1
              WHERE id = $1 AND response_count > 0
              RETURNING response_count`

	var newCount int64
	err := r.db.QueryRow(query, tutorID).Scan(&newCount)
	log.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return errors.New(fmt.Sprintf("no responses or tutor not found %v", tutorID))
		}
		log.Println(err)
		return err
	}

	return nil
}

func (r *Repository) AddResponses(tutorTelegramID int64, responseCount int) (int, error) {
	if responseCount < 1 {
		return 0, errors.New("responses less 0")
	}

	query := `
         	UPDATE tutors
			SET response_count = response_count + $1
			WHERE id = (
				SELECT id FROM users WHERE telegram_id = $2 AND role = $3
			)
			RETURNING response_count`

	var newCount int64
	err := r.db.QueryRow(query, responseCount, tutorTelegramID, "Tutor").Scan(&newCount)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("cannot add responses", tutorTelegramID)
			return 0, errors.New(fmt.Sprintf("tutor with telegram_id %d not found", tutorTelegramID))
		}
		log.Println("cannot add responses", tutorTelegramID)
		return 0, err
	}

	return int(newCount), nil
}

func (r *Repository) UpdateTutorName(tutorID string, name string) error {
	queryUpdateBioTutor := `UPDATE users SET name = $1 WHERE id = $2 AND role = $3`

	_, err := r.db.Exec(queryUpdateBioTutor, name, tutorID, "Tutor")

	if err != nil {
		log.Println(err)
		return errors.New("error update name") //custom_errors.ErrorUpdateBio
	}

	return nil
}
