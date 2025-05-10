package order_testing

import (
	"context"
	config2 "github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"os"
	"testing"
	"time"
)

var configTest = config2.Config{
	DBConfig: config2.DBConfig{
		DBHost:     "localhost",
		DBPort:     "",
		DBUser:     "user",
		DBPassword: "password",
		DBName:     "test",
	},
}

var repo *repository.Repository

func TestMain(m *testing.M) {
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithDatabase(configTest.DBConfig.DBName),
		postgres.WithUsername(configTest.DBConfig.DBUser),
		postgres.WithPassword(configTest.DBConfig.DBPassword),
	)

	if err != nil {
		log.Fatalf("failed start testcontainers: %v", err)
	}

	mappedPort, err := postgresContainer.MappedPort(context.Background(), "5432/tcp")

	if err != nil {
		log.Fatalf("failed mapped port: %v", err)
	}

	configTest.DBConfig.DBPort = mappedPort.Port()

	time.Sleep(10 * time.Second)

	repo = repository.NewRepository(configTest.DBConfig)

	db := repo.GetDB()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			student_id UUID NOT NULL,
			title TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			grade TEXT NOT NULL,
			min_price INTEGER NOT NULL,
			max_price INTEGER NOT NULL,
			tags TEXT[],
			status TEXT NOT NULL,
			response_count INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
	);`)

	if err != nil {
		log.Fatalf("failed init db: %v", err)
	}

	os.Exit(m.Run())
}

func TestRepositoryCreate(t *testing.T) {
	t.Run("Testing create functions", func(t *testing.T) {
		mockOrder := &models.CreateOrder{
			StudentID: "67acb220-7812-4d54-a660-a809b125d088",
			Order: &models.NewOrder{
				Title:       "Help...",
				Name:        "testing exam!",
				Description: "i need test myself:)",
				Grade:       "11",
				MinPrice:    500,
				MaxPrice:    1000,
				Tags:        []string{"math"},
			},
		}

		orderId, err := repo.CreateOrder(mockOrder)

		if err != nil {
			assert.Error(t, err)
		}

		if orderId == "" {
			t.Errorf("Order null ID!")
		}

		order, err := repo.GetOrderByID(orderId)

		if err != nil {
			assert.Error(t, err)
		}

		if orderId != order.ID {
			t.Errorf("OrderID: %v != %v.", orderId, order.ID)
		}
	})
}

func TestRepositoryDelete(t *testing.T) {
	t.Run("Testing delete functions", func(t *testing.T) {
		mockOrder := &models.CreateOrder{
			StudentID: "67acb220-7812-4d54-a660-a809b125d088",
			Order: &models.NewOrder{
				Title:       "Help...",
				Name:        "testing exam!",
				Description: "i need test myself:)",
				Grade:       "11",
				MinPrice:    500,
				MaxPrice:    1000,
				Tags:        []string{"math"},
			},
		}

		orderId, err := repo.CreateOrder(mockOrder)

		if err != nil {
			assert.Error(t, err)
		}

		if orderId == "" {
			t.Errorf("Order null ID!")
		}

		order, err := repo.GetOrderByID(orderId)

		if err != nil {
			assert.Error(t, err)
		}

		if orderId != order.ID {
			t.Errorf("OrderID: %v != %v.", orderId, order.ID)
		}

		err = repo.DeleteOrder(orderId)

		if err != nil {
			assert.Error(t, err)
		}

		order, err = repo.GetOrderByID(orderId)

		if err == nil {
			assert.Error(t, err)
		}
	})
}

func TestRepositoryUpdate(t *testing.T) {
	t.Run("Testing update functions", func(t *testing.T) {
		mockOrder := &models.CreateOrder{
			StudentID: "67acb220-7812-4d54-a660-a809b125d088",
			Order: &models.NewOrder{
				Title:       "Help...",
				Name:        "testing exam!",
				Description: "i need test myself:)",
				Grade:       "11",
				MinPrice:    500,
				MaxPrice:    1000,
				Tags:        []string{"math"},
			},
		}

		orderId, err := repo.CreateOrder(mockOrder)

		if err != nil {
			assert.Error(t, err)
		}

		if orderId == "" {
			t.Errorf("Order null ID!")
		}

		err = repo.UpdateOrder(orderId, &models.UpdateOrder{
			Title: "new_title",
		})

		if err != nil {
			assert.Error(t, err)
		}

		order, err := repo.GetOrderByID(orderId)

		if err != nil {
			assert.Error(t, err)
		}

		if order.Title != mockOrder.Order.Title {
			t.Errorf("Order Title error: %v != %v.", order.Title, mockOrder.Order.Title)
		}
	})
}
