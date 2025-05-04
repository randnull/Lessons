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
	logger.Debug("GetTutorsResponses called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	response, err := r.ResponseService.GetTutorsResponses(UserData)

	if err != nil {
		logger.Error("GetTutorsResponses failed: " + err.Error())
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Debug("GetTutorsResponses successful")

	return ctx.JSON(response)
}

func (r *ResponseController) GetResponseById(ctx *fiber.Ctx) error {
	logger.Debug("GetResponseById called")

	ResponseID := ctx.Params("id")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	response, err := r.ResponseService.GetResponseById(ResponseID, UserData)

	if err != nil {
		logger.Error("GetResponseById failed: " + err.Error())
		if errors.Is(err, custom_errors.ErrorNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Response with id: " + ResponseID + " Not found"})
		} else if errors.Is(err, custom_errors.ErrorServiceError) {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		} else {
			logger.Error("GetResponseById unknown error failed: " + err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unknown error happends"})
		}
	}

	logger.Debug("GetResponseById successful")

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (r *ResponseController) ResponseToOrder(ctx *fiber.Ctx) error {
	logger.Debug("ResponseToOrder called")

	orderID := ctx.Params("id")

	var NewResponse models.NewResponseModel

	if err := ctx.BodyParser(&NewResponse); err != nil {
		logger.Error("ResponseToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	responseID, err := r.ResponseService.ResponseToOrder(orderID, &NewResponse, UserData)

	if err != nil {
		logger.Error("ResponseToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Debug("ResponseToOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"response_id": responseID,
	})
}
