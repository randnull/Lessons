package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
)

type ResponseController struct {
	ResponseService service.ResponseServiceInt
}

func NewResponseController(ResponseServ service.ResponseServiceInt) *ResponseController {
	return &ResponseController{
		ResponseService: ResponseServ,
	}
}

func (r *ResponseController) GetResponseById(ctx *fiber.Ctx) error {
	ResponseID := ctx.Params("id")

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	response, err := r.ResponseService.GetResponseById(ResponseID, UserData)

	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(response)
}

func (r *ResponseController) ResponseToOrder(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	var NewResponse models.NewResponseModel

	if err := ctx.BodyParser(&NewResponse); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	responseID, err := r.ResponseService.ResponseToOrder(orderID, &NewResponse, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"response_id": responseID,
	})
}
