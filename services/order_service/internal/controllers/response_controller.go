package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/logger"
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

func (r *ResponseController) GetTutorsResponses(ctx *fiber.Ctx) error {
	logger.Info("[ResponseController] GetTutorsResponses called")

	UserData, err := getUserData(ctx)
	if err != nil {
		return err
	}

	response, err := r.ResponseService.GetTutorsResponses(UserData)

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("[ResponseController] GetTutorsResponses successful")

	return ctx.JSON(response)
}

func (r *ResponseController) GetResponseById(ctx *fiber.Ctx) error {
	logger.Info("[ResponseController] GetResponseById called")

	UserData, err := getUserData(ctx)
	if err != nil {
		return err
	}

	ResponseID := ctx.Params("id")

	response, err := r.ResponseService.GetResponseById(ResponseID, UserData)

	if err != nil {
		if errors.Is(err, custom_errors.ErrorNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Response with id: " + ResponseID + " Not found"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("[ResponseController] GetResponseById successful")

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *ResponseController) ResponseToOrder(ctx *fiber.Ctx) error {
	logger.Info("[ResponseController] ResponseToOrder called")

	UserData, err := getUserData(ctx)
	if err != nil {
		return err
	}

	orderID := ctx.Params("id")

	var NewResponse models.NewResponseModel

	if err := ctx.BodyParser(&NewResponse); err != nil {
		logger.Error("ResponseToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	responseID, err := r.ResponseService.ResponseToOrder(orderID, &NewResponse, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("[ResponseController] ResponseToOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"response_id": responseID,
	})
}
