package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"log"
	"strings"

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

	userService     service.UserServiceInt
	userControllers *controllers.UserController
}

func NewApp(cfg *config.Config) *App {
	ordergRPC := gRPC_client.NewGRPCClient(cfg.GRPCConfig)
	orderRepo := repository.NewRepository(cfg.DBConfig)
	orderBrokerProducer := rabbitmq.NewRabbitMQ(cfg.MQConfig)

	orderService := service.NewOrderService(orderRepo, orderBrokerProducer, ordergRPC)
	orderController := controllers.NewOrderController(orderService)

	responsesService := service.NewResponseService(orderRepo, orderBrokerProducer, ordergRPC)
	responseControllers := controllers.NewResponseController(responsesService)

	usersService := service.NewUSerService(ordergRPC)
	usersControllers := controllers.NewUserController(usersService)

	return &App{
		repository: orderRepo,

		orderService:     orderService,
		orderControllers: orderController,

		responseService:     responsesService,
		responseControllers: responseControllers,

		userService:     usersService,
		userControllers: usersControllers,

		cfg: cfg,
	}
}

func (a *App) Run() {
	router := fiber.New()

	// В случае плохой производительности - расширить
	//fiber.Config{
	//       Prefork:       true,  // включаем предварительное форкование для увеличения производительности на многоядерных процессорах
	//       ServerHeader:  "Fiber", // добавляем заголовок для идентификации сервера
	//       CaseSensitive: true,    // включаем чувствительность к регистру в URL
	//       StrictRouting: true,    // включаем строгую маршрутизацию
	//   })

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*", // НЕБЕЗОПАСНО, ЗАМЕНИТЬ ТОЛЬКО НА ХОСТ ФРОНТА!
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: "*",
	}))

	router.Use(logger.New(logger.Config{
		Format: "[LOG] ${time} [${ip}]:${port} ${status} - ${method} - ${latency} ${path}\n",
		// output ...
	}))

	router.Get("/", a.orderControllers.HealtzHandler)

	orders := router.Group("api/orders")
	orders.Use(controllers.TokenAuthMiddleware(a.cfg.BotConfig))

	orders.Post("/", a.orderControllers.CreateOrder)                      // StudentOnly
	orders.Get("/id/:id", a.orderControllers.GetOrderByID)                // All
	orders.Get("/", a.orderControllers.GetAllUsersOrders)                 // All
	orders.Delete("/id/:id", a.orderControllers.DeleteOrderByID)          // StudentOnly
	orders.Get("/all", a.orderControllers.GetAllOrders)                   // TutorOnly
	orders.Put("/id/:id", a.orderControllers.UpdateOrderByID)             // StudentOnly
	orders.Get("/mini/id/:id/", a.orderControllers.GetOrderByIdTutor)     // Tutor Only
	orders.Post("/select/id/:id/", a.orderControllers.SelectTutorToOrder) // Student Only

	responses := router.Group("api/responses")
	responses.Use(controllers.TokenAuthMiddlewareResponses(a.cfg.BotConfig)) // другой bot config

	responses.Post("/id/:id", a.responseControllers.ResponseToOrder)
	responses.Get("/id/:id", a.responseControllers.GetResponseById)

	users := router.Group("/api/users")
	users.Use(controllers.TokenAuthMiddlewareResponses(a.cfg.BotConfig)) // другой bot config

	//users.Post("/", a.userControllers.CreateUser)
	users.Get("/", a.userControllers.GetAllUser)
	users.Get("/id/:id", a.userControllers.GetUser)
	users.Get("/tutor/id/:id", a.userControllers.GetTutor)
	users.Post("/tutor/bio/id/:id", a.userControllers.UpdateBioTutor)

	ListenPort := fmt.Sprintf(":%v", a.cfg.ServerPort)

	log.Fatal(router.Listen(ListenPort)) // graceful shutdown
}
