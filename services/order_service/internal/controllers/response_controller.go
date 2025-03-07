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
	orderID := ctx.Params("id")

	var NewResponse models.NewResponseModel

	//if err := ctx.BodyParser(&NewResponse); err != nil {
	//	return fiber.NewError(fiber.StatusBadRequest, err.Error())
	//}

	NewResponse.OrderId = orderID

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	responseID, err := r.ResponseService.ResponseToOrder(&NewResponse, InitData)

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"response_id": responseID,
	})
}
