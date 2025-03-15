package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
)

type UserController struct {
	UserService service.UserServiceInt
}

func NewUserController(UserServ service.UserServiceInt) *UserController {
	return &UserController{
		UserService: UserServ,
	}
}

//func (u *UserController) CreateUser(ctx *fiber.Ctx) error {
//	UserData, ok := ctx.Locals("user_data").(models.UserData)
//
//	if !ok {
//		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
//	}
//
//	UserID, err := u.UserService.CreateUser(UserData)
//	if err != nil {
//		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error with create user"})
//	}
//
//	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
//		"userID": UserID,
//	})
//}

func (u *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	user, err := u.UserService.GetUser(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return ctx.JSON(user)
}

func (u *UserController) GetAllUser(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	users, err := u.UserService.GetAllUsers()

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "users not found"})
	}

	return ctx.JSON(users)
}
