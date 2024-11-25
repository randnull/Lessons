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
	var InitDataFromUser models.InitModel

	if err := ctx.BodyParser(&InitDataFromUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad initdata"})
	}

	if InitDataFromUser.InitData == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Field 'initData' is required"})
	}

	fmt.Println("INIT DATA ON HANDLER:", InitDataFromUser)

	jwt_token, err := c.Service.Login(InitDataFromUser.InitData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(models.ResponseToken{Token: jwt_token})
}
