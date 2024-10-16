package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/controllers"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/service"
)

type App struct {
	repository  repository.OrderRepository
	service     service.OrderServiceInt
	controllers *controllers.OrderController
}

func NewApp() *App {
	orderRepo := repository.NewRepository()
	orderService := service.NewOrderService(orderRepo)
	orderController := controllers.NewOrderController(orderService)

	return &App{
		repository:  orderRepo,
		service:     orderService,
		controllers: orderController,
	}
}

func (a *App) Run() {
	router := fiber.New()

	router.Post("/orders", a.controllers.CreateOrder)
	router.Get("/orders/:id", a.controllers.GetOrderByID)

	log.Fatal(router.Listen(":3000"))
}
