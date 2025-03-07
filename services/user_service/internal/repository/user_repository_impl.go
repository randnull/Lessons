package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/custom_errors"
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
	ExistedUser, err := r.GetUserInfoById(user.TelegramId)
	if err == nil {
		return ExistedUser.Id, nil
	}

	query := `INSERT INTO users (telegram_id, name, created_at) VALUES ($1, $2, $3) RETURNING id`

	var UserId string

	currentTime := time.Now()

	err = r.db.QueryRow(query,
		user.TelegramId,
		user.Name,
		currentTime,
	).Scan(&UserId)
	fmt.Println(err)
	if err != nil {
		return "", custom_errors.ErrorWithCreate
	}

	return UserId, nil
}

func (r *Repository) GetUserInfoById(telegramID int64) (*models.UserDB, error) {
	user := &models.UserDB{}

	query := `SELECT id, telegram_id, name, created_at FROM users WHERE telegram_id = $1`

	err := r.db.QueryRow(query, telegramID).Scan(
		&user.Id,
		&user.TelegramID,
		&user.Name,
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
