package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	"log"
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

func (u *UserController) GetMyTutorProfile(ctx *fiber.Ctx) error {
	log.Println("Запрос на всей инфы от репета своей страницы")

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad auth data"})
	}

	tutorID := UserData.UserID

	info, err := u.UserService.GetTutorInfoById(tutorID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get data tutor" + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(info)
}

func (u *UserController) GetTutorInfoById(ctx *fiber.Ctx) error {
	log.Println("Запрос на всей инфы от репета")

	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad auth data"})
	}

	tutorID := ctx.Params("id")

	info, err := u.UserService.GetTutorInfoById(tutorID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get data tutor" + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(info)
}

// Запросы для обнолвения параметр от репетиторов
func (u *UserController) UpdateTagsTutor(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad auth data"})
	}

	log.Println("Запрос на обновление тегов")

	var UpdateTagsTutor models.UpdateTagsTutor

	if err := ctx.BodyParser(&UpdateTagsTutor); err != nil {
		log.Println(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad format"})
	}

	tutorID := UserData.UserID

	success, err := u.UserService.UpdateTagsTutor(UpdateTagsTutor.Tags, tutorID)
	if err != nil || !success {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update tags"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
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

// Отзывы
func (u *UserController) CreateReview(ctx *fiber.Ctx) error {
	log.Println("Запрос на создание отзыва")

	var ReviewRequest models.ReviewRequest

	if err := ctx.BodyParser(&ReviewRequest); err != nil {
		log.Println("Ошибка парсинга:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad format"})
	}

	UserData, ok := ctx.Locals("user_data").(models.UserData)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad auth"})
	}

	id, err := u.UserService.CreateReview(UserData.UserID, ReviewRequest.TutorID, ReviewRequest.Comment, ReviewRequest.Rating)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error review" + err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (u *UserController) GetReviewsByTutor(ctx *fiber.Ctx) error {
	log.Println("Запрос на получения отзывов")
	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad auth data"})
	}

	tutorID := ctx.Params("tutor_id")

	reviews, err := u.UserService.GetReviewsByTutor(tutorID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get auth" + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(reviews)
}

func (u *UserController) GetReviewByID(ctx *fiber.Ctx) error {
	log.Println("Запрос на получения отзыва id")

	_, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad auth data"})
	}

	reviewID := ctx.Params("id")

	review, err := u.UserService.GetReviewsByID(reviewID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get review" + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(review)
}

func (u *UserController) ChangeTutorActive(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	var IsActive models.ChangeActive

	if err := ctx.BodyParser(&IsActive); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	IsOk, err := u.UserService.ChangeTutorActive(IsActive.IsActive, UserData)

	if err != nil || !IsOk {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Cannot update status"})
	}

	return nil
}

func (u *UserController) UpdateNameTutor(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	var UpdateNameTutor models.UpdateNameTutor

	if err := ctx.BodyParser(&UpdateNameTutor); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := u.UserService.UpdateTutorName(UserData.UserID, UpdateNameTutor.Name)

	if err != nil {
		return err
	}

	return nil
}
