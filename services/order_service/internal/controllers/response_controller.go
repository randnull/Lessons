package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type ResponseController struct {
	ResponseService service.ResponseServiceInt
}

func NewResponseController(ResponseServ service.ResponseServiceInt) *ResponseController {
	return &ResponseController{
		ResponseService: ResponseServ,
	}
}

func (r *ResponseController) ResponseToOrder(ctx *fiber.Ctx) error {
	var NewResponse models.NewResponseModel

	if err := ctx.BodyParser(&NewResponse); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	err := r.ResponseService.ResponseToOrder(&NewResponse, InitData)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
