package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	"github.com/randnull/Lessons/internal/utils"
	"github.com/randnull/Lessons/pkg/logger"
)

type ResponseController struct {
	ResponseService service.ResponseServiceInt
}

func NewResponseController(ResponseServ service.ResponseServiceInt) *ResponseController {
	return &ResponseController{
		ResponseService: ResponseServ,
	}
}

func (r *ResponseController) GetTutorsResponses(ctx *fiber.Ctx) error {
	logger.Info("[ResponseController] GetTutorsResponses called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	response, err := r.ResponseService.GetTutorsResponses(UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[ResponseController] GetTutorsResponses successful")

	return ctx.JSON(response)
}

func (r *ResponseController) GetResponseById(ctx *fiber.Ctx) error {
	logger.Info("[ResponseController] GetResponseById called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	ResponseID := ctx.Params("id")

	response, err := r.ResponseService.GetResponseById(ResponseID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[ResponseController] GetResponseById successful")

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *ResponseController) ResponseToOrder(ctx *fiber.Ctx) error {
	logger.Info("[ResponseController] ResponseToOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orderID := ctx.Params("id")

	var NewResponse models.NewResponseModel

	if err := ctx.BodyParser(&NewResponse); err != nil {
		logger.Error("ResponseToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	responseID, err := r.ResponseService.ResponseToOrder(orderID, &NewResponse, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[ResponseController] ResponseToOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"response_id": responseID,
	})
}
