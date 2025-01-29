package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"log"

	"github.com/randnull/Lessons/internal/controllers"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/service"
)

type App struct {
	cfg        *config.Config
	repository repository.OrderRepository

	orderService     service.OrderServiceInt
	orderControllers *controllers.OrderController

	responseService     service.ResponseServiceInt
	responseControllers *controllers.ResponseController
}

func NewApp(cfg *config.Config) *App {
	orderRepo := repository.NewRepository(cfg.DBConfig)
	orderBrokerProducer := rabbitmq.NewRabbitMQ(cfg.MQConfig)

	orderService := service.NewOrderService(orderRepo, orderBrokerProducer)
	orderController := controllers.NewOrderController(orderService)

	responsesService := service.NewResponseService(orderRepo, orderBrokerProducer)
	responseControllers := controllers.NewResponseController(responsesService)

	return &App{
		repository: orderRepo,

		orderService:     orderService,
		orderControllers: orderController,

		responseService:     responsesService,
		responseControllers: responseControllers,

		cfg: cfg,
	}
}

func (a *App) Run() {
	router := fiber.New()

	router.Use(logger.New())

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*", // НЕБЕЗОПАСНО, ЗАМЕНИТЬ ТОЛЬКО НА ХОСТ ФРОНТА!
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "*",
	}))

	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	router.Get("/", a.orderControllers.HealtzHandler)

	orders := router.Group("/orders")
	orders.Use(controllers.TokenAuthMiddleware(a.cfg.BotConfig))

	orders.Post("/", a.orderControllers.CreateOrder)
	orders.Get("/id/:id", a.orderControllers.GetOrderByID)
	orders.Get("/", a.orderControllers.GetAllUsersOrders)
	orders.Delete("/id/:id", a.orderControllers.DeleteOrderByID)

	orders.Get("/all", a.orderControllers.GetAllOrders)

	// Put or patch

	responses := router.Group("/responses")
	responses.Use(controllers.TokenAuthMiddlewareResponses(a.cfg.BotConfig)) // другой bot config

	responses.Post("/id/:id", a.responseControllers.ResponseToOrder)

	ListenPort := fmt.Sprintf(":%v", a.cfg.ServerPort)

	log.Fatal(router.Listen(ListenPort)) // graceful shutdown
}
