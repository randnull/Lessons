package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/controllers"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/service"
	"log"
)

type App struct {
	cfg         *config.Config
	service     service.AuthServiceInt
	controllers *controllers.AuthHandlers
}

func NewApp(cfg *config.Config) *App {
	gRPCClient := gRPC_client.NewGRPCClient(cfg.GRPCConfig)

	authService := service.NewAuthService(&cfg.JWTConfig, gRPCClient)
	authControllers := controllers.NewUserHandler(authService, cfg)

	auth_app := &App{
		cfg:         cfg,
		service:     authService,
		controllers: authControllers,
	}

	return auth_app
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

	router.Post("/auth/init-data", a.controllers.Login)

	addr := fmt.Sprintf(":%v", a.cfg.ServerPort)

	log.Printf("Listen on: %s\n", addr)

	log.Fatal(router.Listen(addr))
}
