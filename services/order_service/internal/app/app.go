package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	lg "github.com/randnull/Lessons/pkg/logger"
	"github.com/randnull/Lessons/pkg/rabbitmq"
	"time"

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

	usersService := service.NewUserService(ordergRPC, orderBrokerProducer, orderRepo)
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

func (a *App) Run(ctx context.Context) error {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "*",
	}))

	router.Use(logger.New(logger.Config{
		Format: "[LOG] ${time} [${ip}]:${port} ${status} - ${method} - ${latency} ${path}\n",
	}))

	orders := router.Group("api/orders")

	studentType := models.StudentType
	tutorType := models.TutorType
	adminType := models.AdminType
	anyType := models.AnyType

	orders.Post("/", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.CreateOrder)
	//orders.Get("/", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.GetAllUsersOrders)
	orders.Get("/pagination/student", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.GetStudentOrdersPagination)
	orders.Put("/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.UpdateOrderByID)
	orders.Get("/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.GetOrderByID)
	orders.Delete("/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.DeleteOrderByID)
	orders.Post("/id/:id/active", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.SetActiveToOrder)
	orders.Post("/select/id/:id/", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.SelectTutorToOrder)
	orders.Post("/suggest/tutor/:id/", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.orderControllers.SuggestOrder)

	orders.Get("/tutor/id/:id/", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.orderControllers.GetOrderByIdTutor)
	orders.Get("/pagination", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.orderControllers.GetOrdersPagination)

	responses := router.Group("api/responses")

	responses.Post("/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.responseControllers.ResponseToOrder)
	responses.Get("/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, anyType), a.responseControllers.GetResponseById)
	responses.Get("/list", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.responseControllers.GetTutorsResponses)

	users := router.Group("/api/users")

	users.Get("/pagination", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.userControllers.GetTutorsPagination)

	users.Get("/tutor/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, anyType), a.userControllers.GetTutorInfoById)
	users.Get("/tutor/profile", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.userControllers.GetMyTutorProfile)

	users.Post("/tutor/bio", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.userControllers.UpdateBioTutor)
	users.Post("/tutor/tags", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.userControllers.UpdateTagsTutor)
	users.Post("/tutor/active", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.userControllers.ChangeTutorActive)
	users.Post("/tutor/name", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.userControllers.UpdateNameTutor)

	users.Post("/review", controllers.TokenAuthMiddleware(a.cfg.BotConfig, studentType), a.userControllers.CreateReview)
	users.Get("/review/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, anyType), a.userControllers.GetReviewByID)
	users.Get("/tutor/id/:id/reviews", controllers.TokenAuthMiddleware(a.cfg.BotConfig, anyType), a.userControllers.GetReviewsByTutor)
	users.Post("/review/activate", controllers.TokenAuthMiddleware(a.cfg.BotConfig, tutorType), a.userControllers.SetReviewActive)

	admins := router.Group("api/admins")

	admins.Post("/ban/user/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.userControllers.BanUser)
	admins.Post("/ban/order/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.orderControllers.SetBanOrder)
	admins.Post("/approve/order/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.orderControllers.SetApprovedOrder)
	admins.Get("/orders", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.orderControllers.GetAllOrders)
	admins.Get("/order/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.orderControllers.GetOrderByIdTutor)
	admins.Get("/users", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.userControllers.GetAllUsers)
	admins.Get("/user/id/:id", controllers.TokenAuthMiddleware(a.cfg.BotConfig, adminType), a.userControllers.GetUserById)

	ListenPort := fmt.Sprintf(":%v", a.cfg.ServerPort)

	lg.Error(fmt.Sprintf("server starting listening port: %v", ListenPort))

	go func() {
		err := router.Listen(ListenPort)
		if err != nil {
			lg.Error(fmt.Sprintf("server stopped with error: %v", err))
		}
	}()

	<-ctx.Done()

	lg.Info("server graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := router.ShutdownWithContext(shutdownCtx)
	if err != nil {
		lg.Info(fmt.Sprintf("server shutdown error: %v", err))
		return err
	}

	return nil
}
