package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	"strconv"
)

type UserController struct {
	UserService service.UserServiceInt
}

func NewUserController(UserServ service.UserServiceInt) *UserController {
	return &UserController{
		UserService: UserServ,
	}
}

func (u *UserController) GetUser(ctx *fiber.Ctx) error {
	// ЗАКРЫТЬ ЭТО ДЛЯ ВНЕШКИ
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

func (u *UserController) GetTutor(ctx *fiber.Ctx) error {
	// ЗАКРЫТЬ ЭТО ДЛЯ ВНЕШКИ
	id := ctx.Params("id")

	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	user, err := u.UserService.GetTutor(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tutor not found"})
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

func (u *UserController) UpdateBioTutor(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	var UpdateBioModel models.UpdateBioTutor

	if err := ctx.BodyParser(&UpdateBioModel); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := u.UserService.UpdateBioTutor(UpdateBioModel, UserData)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserController) GetTutorsPagination(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})

	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})

	}

	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	users, err := u.UserService.GetAllTutorsPagination(page, size)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "users not found"})
	}

	return ctx.JSON(users)
}
