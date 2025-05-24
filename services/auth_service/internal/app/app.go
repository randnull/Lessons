package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/controllers"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/service"
	lg "github.com/randnull/Lessons/pkg/logger"
	"log"
	"time"
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

	authApp := &App{
		cfg:         cfg,
		service:     authService,
		controllers: authControllers,
	}

	return authApp
}

func (a *App) Run(ctx context.Context) error {
	router := fiber.New()

	router.Use(logger.New())

	corsOriginService := a.cfg.CorsOrigin

	router.Use(cors.New(cors.Config{
		AllowOrigins: corsOriginService,
		AllowMethods: "GET,POST",
		AllowHeaders: "*",
	}))

	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	router.Post("/auth/init-data", a.controllers.Login)

	ListenPort := fmt.Sprintf(":%v", a.cfg.ServerPort)

	go func() {
		err := router.Listen(ListenPort)
		if err != nil {
			log.Printf("server stopped with error: %v", err)
		}
	}()

	<-ctx.Done()

	lg.Info("server graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := router.ShutdownWithContext(shutdownCtx)
	if err != nil {
		log.Printf("server shutdown error: %v", err)
		return err
	}

	return nil
}
