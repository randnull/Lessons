package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	"github.com/randnull/Lessons/internal/utils"
	"github.com/randnull/Lessons/pkg/logger"
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

func (u *UserController) GetUserById(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetUserById called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	userID := ctx.Params("id")

	user, err := u.UserService.GetUserById(userID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetUserById successful")

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (u *UserController) GetTutorsPagination(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetTutorsPagination called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})
	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	tag := ctx.Query("tag")

	users, err := u.UserService.GetAllTutorsPagination(page, size, tag, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetTutorsPagination successful")

	return ctx.JSON(users)
}

func (u *UserController) GetMyTutorProfile(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetMyTutorProfile called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	tutorID := UserData.UserID

	info, err := u.UserService.GetTutorInfoById(tutorID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetMyTutorProfile successful")

	return ctx.Status(fiber.StatusOK).JSON(info)
}

func (u *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetAllUsers called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	users, err := u.UserService.GetAllUsers(UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetAllUsers successful")

	return ctx.Status(fiber.StatusOK).JSON(users)
}

func (u *UserController) GetTutorInfoById(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetTutorInfoById called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	tutorID := ctx.Params("id")

	info, err := u.UserService.GetTutorInfoById(tutorID, UserData)

	info.Tutor.TelegramID = 0

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetTutorInfoById successful")

	return ctx.Status(fiber.StatusOK).JSON(info)
}

func (u *UserController) BanUser(ctx *fiber.Ctx) error {
	logger.Info("[UserController] BanUser called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var BanUser models.BanUser

	if err := ctx.BodyParser(&BanUser); err != nil {
		logger.Error("[UserController] BanUser failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err = u.UserService.BanUser(BanUser, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] BanUser successful")

	return ctx.SendStatus(fiber.StatusOK)
}

func (u *UserController) UpdateTagsTutor(ctx *fiber.Ctx) error {
	logger.Info("[UserController] UpdateTagsTutor called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var UpdateTagsTutor models.UpdateTagsTutor

	if err := ctx.BodyParser(&UpdateTagsTutor); err != nil {
		logger.Error("[UserController] UpdateTagsTutor failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	if err := models.Valid.Struct(UpdateTagsTutor); err != nil {
		logger.Error("[UserController] UpdateTagsTutor failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = u.UserService.UpdateTagsTutor(UpdateTagsTutor.Tags, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] UpdateTagsTutor successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

func (u *UserController) UpdateBioTutor(ctx *fiber.Ctx) error {
	logger.Info("[UserController] UpdateBioTutor called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var UpdateBioModel models.UpdateBioTutor

	if err := ctx.BodyParser(&UpdateBioModel); err != nil {
		logger.Error("[UserController] UpdateBioTutor failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := models.Valid.Struct(UpdateBioModel); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = u.UserService.UpdateBioTutor(UpdateBioModel, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] UpdateBioTutor successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) UpdateNameTutor(ctx *fiber.Ctx) error {
	logger.Info("[UserController] UpdateNameTutor called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var UpdateNameTutor models.UpdateNameTutor

	if err := ctx.BodyParser(&UpdateNameTutor); err != nil {
		logger.Error("[UserController] UpdateNameTutor failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := models.Valid.Struct(UpdateNameTutor); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = u.UserService.UpdateTutorName(UserData.UserID, UpdateNameTutor.Name, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] UpdateNameTutor successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) ChangeTutorActive(ctx *fiber.Ctx) error {
	logger.Info("[UserController] ChangeTutorActive called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var IsActive models.ChangeActive

	if err := ctx.BodyParser(&IsActive); err != nil {
		logger.Error("[UserController] ChangeTutorActive failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	if err := models.Valid.Struct(IsActive); err != nil {
		logger.Error("[UserController] ChangeTutorActive failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = u.UserService.ChangeTutorActive(IsActive.IsActive, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] ChangeTutorActive successful")

	return ctx.SendStatus(fiber.StatusCreated)
}

func (u *UserController) GetReviewsByTutor(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetReviewsByTutor called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	tutorID := ctx.Params("tutor_id")

	reviews, err := u.UserService.GetReviewsByTutor(tutorID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetReviewsByTutor successful")

	return ctx.Status(fiber.StatusOK).JSON(reviews)
}

func (u *UserController) CreateReview(ctx *fiber.Ctx) error {
	logger.Info("[UserController] CreateReview called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var ReviewRequest models.ReviewRequest

	if err := ctx.BodyParser(&ReviewRequest); err != nil {
		logger.Error("CreateReview failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := models.Valid.Struct(ReviewRequest); err != nil {
		logger.Error("CreateReview failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := u.UserService.CreateReview(ReviewRequest, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] CreateReview successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (u *UserController) GetReviewByID(ctx *fiber.Ctx) error {
	logger.Info("[UserController] GetReviewByID called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	reviewID := ctx.Params("id")

	review, err := u.UserService.GetReviewsByID(reviewID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] GetReviewByID successful")

	return ctx.Status(fiber.StatusOK).JSON(review)
}

func (u *UserController) SetReviewActive(ctx *fiber.Ctx) error {
	logger.Info("[UserController] SetReviewActive called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var ReviewActive models.ReviewActive

	if err := ctx.BodyParser(&ReviewActive); err != nil {
		logger.Error("[UserController] SetReviewActive parse failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := models.Valid.Struct(ReviewActive); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = u.UserService.SetReviewActive(ReviewActive.ReviewID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	logger.Info("[UserController] SetReviewActive successful")

	return ctx.SendStatus(fiber.StatusCreated)
}
