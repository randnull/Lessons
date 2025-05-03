package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/logger"
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

// Запросы про репетиторов
func (u *UserController) GetTutorsPagination(ctx *fiber.Ctx) error {
	logger.Debug("GetTutorsPagination called")

	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		logger.Error("GetTutorsPagination failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})
	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		logger.Error("GetTutorsPagination failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	tag := ctx.Query("tag")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	users, err := u.UserService.GetAllTutorsPagination(page, size, tag, UserData)

	if err != nil {
		logger.Error("GetTutorsPagination failed: " + err.Error())

		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "users not found"})
	}

	logger.Debug("GetTutorsPagination successful")

	return ctx.JSON(users)
}

func (u *UserController) GetMyTutorProfile(ctx *fiber.Ctx) error {
	logger.Debug("GetMyTutorProfile called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	tutorID := UserData.UserID

	info, err := u.UserService.GetTutorInfoById(tutorID, UserData)
	if err != nil {
		logger.Error("GetMyTutorProfile failed: " + err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get data tutor" + err.Error()})
	}

	logger.Debug("GetMyTutorProfile successful")

	return ctx.Status(fiber.StatusOK).JSON(info)
}

func (u *UserController) GetTutorInfoById(ctx *fiber.Ctx) error {
	logger.Debug("GetTutorInfoById called")

	tutorID := ctx.Params("id")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	info, err := u.UserService.GetTutorInfoById(tutorID, UserData)

	info.Tutor.TelegramID = 0

	if err != nil {
		logger.Error("GetTutorInfoById failed: " + err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get data tutor" + err.Error()})
	}

	logger.Debug("GetTutorInfoById successful")

	return ctx.Status(fiber.StatusOK).JSON(info)
}

func (u *UserController) UpdateTagsTutor(ctx *fiber.Ctx) error {
	logger.Debug("UpdateTagsTutor called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	var UpdateTagsTutor models.UpdateTagsTutor

	if err := ctx.BodyParser(&UpdateTagsTutor); err != nil {
		logger.Error("UpdateTagsTutor failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad format"})
	}

	success, err := u.UserService.UpdateTagsTutor(UpdateTagsTutor.Tags, UserData)
	if err != nil || !success {
		logger.Error("UpdateTagsTutor failed: " + err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update tags"})
	}

	logger.Debug("UpdateTagsTutor successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

func (u *UserController) UpdateBioTutor(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Debug("UpdateBioTutor called")

	var UpdateBioModel models.UpdateBioTutor

	if err := ctx.BodyParser(&UpdateBioModel); err != nil {
		logger.Error("UpdateBioTutor failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := u.UserService.UpdateBioTutor(UpdateBioModel, UserData)

	if err != nil {
		return err
	}
	logger.Debug("UpdateBioTutor successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

// Отзывы
func (u *UserController) CreateReview(ctx *fiber.Ctx) error {
	logger.Debug("CreateReview called")

	var ReviewRequest models.ReviewRequest

	if err := ctx.BodyParser(&ReviewRequest); err != nil {
		logger.Error("CreateReview failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad format"})
	}

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	id, err := u.UserService.CreateReview(ReviewRequest, UserData)
	if err != nil {
		logger.Error("CreateReview failed: " + err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error review" + err.Error()})
	}

	logger.Debug("CreateReview successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (u *UserController) GetReviewsByTutor(ctx *fiber.Ctx) error {
	logger.Debug("GetReviewsByTutor called")

	tutorID := ctx.Params("tutor_id")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	reviews, err := u.UserService.GetReviewsByTutor(tutorID, UserData)

	if err != nil {
		logger.Error("GetReviewsByTutor failed: " + err.Error())

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get auth" + err.Error()})
	}
	logger.Debug("GetReviewsByTutor successful")

	return ctx.Status(fiber.StatusOK).JSON(reviews)
}

func (u *UserController) GetReviewByID(ctx *fiber.Ctx) error {
	logger.Debug("GetReviewByID called")

	reviewID := ctx.Params("id")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	review, err := u.UserService.GetReviewsByID(reviewID, UserData)

	if err != nil {
		logger.Error("GetReviewByID failed: " + err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get review" + err.Error()})
	}
	logger.Debug("GetReviewByID successful")

	return ctx.Status(fiber.StatusOK).JSON(review)
}

func (u *UserController) ChangeTutorActive(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Debug("ChangeTutorActive called")

	var IsActive models.ChangeActive

	if err := ctx.BodyParser(&IsActive); err != nil {
		logger.Error("ChangeTutorActive failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	IsOk, err := u.UserService.ChangeTutorActive(IsActive.IsActive, UserData)

	if err != nil || !IsOk {
		logger.Error("ChangeTutorActive failed: " + err.Error())
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Cannot update status"})
	}

	logger.Debug("ChangeTutorActive successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) UpdateNameTutor(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Debug("UpdateNameTutor called")

	var UpdateNameTutor models.UpdateNameTutor

	if err := ctx.BodyParser(&UpdateNameTutor); err != nil {
		logger.Error("UpdateNameTutor failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := u.UserService.UpdateTutorName(UserData.UserID, UpdateNameTutor.Name, UserData)

	if err != nil {
		return err
	}

	logger.Debug("UpdateNameTutor successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) SetReviewActive(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Debug("SetReviewActive called")

	var ReviewActive models.ReviewActive

	if err := ctx.BodyParser(&ReviewActive); err != nil {
		logger.Error("SetReviewActive parse failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := u.UserService.SetReviewActive(ReviewActive.ReviewID, UserData)

	if err != nil {
		logger.Error("SetReviewActive failed: " + err.Error())
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Cannot set review active status"})
	}

	logger.Debug("SetReviewActive successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) BanUser(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Debug("BanUser called")

	var BanUser models.BanUser

	if err := ctx.BodyParser(&BanUser); err != nil {
		logger.Error("BanUser parse failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := u.UserService.BanUser(BanUser, UserData)

	if err != nil {
		logger.Error("BanUser failed: " + err.Error())
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Cannot set review active status"})
	}

	logger.Debug("BanUser successful")

	return ctx.SendStatus(fiber.StatusOK)
}
