package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/randnull/Lessons/internal/config"
	"log"

	"github.com/randnull/Lessons/internal/controllers"
	"github.com/randnull/Lessons/internal/repository"
	"github.com/randnull/Lessons/internal/service"
)

type App struct {
	cfg         *config.Config
	repository  repository.OrderRepository
	service     service.OrderServiceInt
	controllers *controllers.OrderController
}

func NewApp(cfg *config.Config) *App {
	orderRepo := repository.NewRepository(cfg.DBConfig)
	orderService := service.NewOrderService(orderRepo)
	orderController := controllers.NewOrderController(orderService)

	return &App{
		repository:  orderRepo,
		service:     orderService,
		controllers: orderController,
		cfg:         cfg,
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

	router.Get("/", a.controllers.HealtzHandler)

	orders := router.Group("/orders")
	orders.Use(controllers.TokenAuthMiddleware(a.cfg.BotConfig))

	orders.Post("/", a.controllers.CreateOrder)
	orders.Get("/:id", a.controllers.GetOrderByID)
	orders.Get("/", a.controllers.GetAllUsersOrders)
	orders.Delete("/:id", a.controllers.DeleteOrderByID)
	// Put or patch

	log.Fatal(router.Listen(":8001"))
}
