package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/randnull/Lessons/internal/models"
	"log"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository() *Repository {
	//link := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
	//	"CHANGE", "CHANGE", "postgresql", "5432", "user_database")
	//
	//db, err := sqlx.Open("postgres", link)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = db.PingContext(context.Background())
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Print("Database is ready")

	return &Repository{
		db: nil,
	}
}

//func (r *Repository) CreateUser(user *models.User) (string, error) {
//	return "", nil
//}

func (r *Repository) GetUserById(user_id string) (*models.User, error) {
	return &models.User{
		UserId: "23",
		Name:   "qe",
	}, nil
}
