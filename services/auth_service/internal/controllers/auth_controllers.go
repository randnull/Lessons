package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/config"
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

	if err := ctx.BodyParser(&AuthData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad initdata"})
	}

	if AuthData.InitData == "" || AuthData.Role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Field 'initData' and 'role' is required"})
	}

	fmt.Println("INIT DATA ON HANDLER:", AuthData)

	jwtToken, err := c.Service.Login(&AuthData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": jwtToken,
	})
}
