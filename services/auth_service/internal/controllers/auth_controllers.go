package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/config"
	lg "github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
)

type AuthHandlers struct {
	Service service.AuthServiceInt
	cfg     *config.Config
}

func NewUserHandler(AuthServ service.AuthServiceInt, cfg *config.Config) *AuthHandlers {
	return &AuthHandlers{
		Service: AuthServ,
		cfg:     cfg,
	}
}

func (c *AuthHandlers) Login(ctx *fiber.Ctx) error {
	var AuthData models.AuthData

	lg.Info("Login called")

	if err := ctx.BodyParser(&AuthData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad initdata"})
	}

	if AuthData.InitData == "" || AuthData.Role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Field initData and role is empty"})
	}

	jwtToken, err := c.Service.Login(&AuthData)

	if err != nil {
		lg.Error(fmt.Sprintf("Error auth. Error: %v", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error with auth"})
	}

	lg.Info("jwt created")
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": jwtToken,
	})
}
