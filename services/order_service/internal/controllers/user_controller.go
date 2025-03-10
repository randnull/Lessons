package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/service"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type UserController struct {
	UserService service.UserServiceInt
}

func NewUserController(UserServ service.UserServiceInt) *UserController {
	return &UserController{
		UserService: UserServ,
	}
}

func (u *UserController) CreateUser(ctx *fiber.Ctx) error {
	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	UserID, err := u.UserService.CreateUser(InitData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error with create user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"userID": UserID,
	})
}

func (u *UserController) GetUser(ctx *fiber.Ctx) error {
	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	user, err := u.UserService.GetUser(InitData.User.ID)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return ctx.JSON(user)
}
