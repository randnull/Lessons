package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	"github.com/randnull/Lessons/internal/utils"
	custom_logger "github.com/randnull/Lessons/pkg/logger"
	"strconv"
)

type OrderController struct {
	OrderService service.OrderServiceInt
}

func NewOrderController(OrderService service.OrderServiceInt) *OrderController {
	return &OrderController{
		OrderService: OrderService,
	}
}

func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] CreateOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	var order models.NewOrder

	if err := ctx.BodyParser(&order); err != nil {
		custom_logger.Error("[OrderController] CreateOrder parse failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	if err := models.Valid.Struct(order); err != nil {
		custom_logger.Error("[OrderController] CreateOrder valid failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	orderID, err := c.OrderService.CreateOrder(&order, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] CreateOrder successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"orderID": orderID,
	})
}

func (c *OrderController) GetOrderByIdTutor(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] GetOrderByIdTutor called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	id := ctx.Params("id")

	order, err := c.OrderService.GetOrderByIdTutor(id, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] GetOrderByIdTutor successful")

	return ctx.Status(fiber.StatusOK).JSON(order)
}

func (c *OrderController) GetOrderByID(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] GetOrderByIdTutor called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	id := ctx.Params("id")

	order, err := c.OrderService.GetOrderById(id, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] GetOrderByID successful")

	return ctx.Status(fiber.StatusOK).JSON(order)
}

func (c *OrderController) GetOrdersPagination(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] GetOrdersPagination called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "page param not correct"})
	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	if size > 100 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param more than maximum"})
	}

	tag := ctx.Query("tag")

	orders, err := c.OrderService.GetOrdersWithPagination(page, size, tag, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
	}
	custom_logger.Info("[OrderController] GetOrdersPagination successful")

	return ctx.Status(fiber.StatusOK).JSON(orders)
}

func (c *OrderController) GetStudentOrdersPagination(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] GetStudentOrdersPagination called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "page param not correct"})
	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	if size > 100 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param more than maximum"})
	}

	orders, err := c.OrderService.GetStudentOrdersWithPagination(page, size, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] GetStudentOrdersPagination successful")

	return ctx.Status(fiber.StatusOK).JSON(orders)
}

func (c *OrderController) DeleteOrderByID(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] DeleteOrderByID called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orderID := ctx.Params("id")

	err = c.OrderService.DeleteOrder(orderID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] DeleteOrderByID successful")

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (c *OrderController) UpdateOrderByID(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] UpdateOrderByID called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orderID := ctx.Params("id")

	var order models.UpdateOrder

	if err := ctx.BodyParser(&order); err != nil {
		custom_logger.Error("[OrderController] UpdateOrderByID failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = c.OrderService.UpdateOrder(orderID, &order, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] UpdateOrderByID successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SelectTutorToOrder(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] SelectTutorToOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	responseID := ctx.Params("id")

	err = c.OrderService.SelectTutor(responseID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] SelectTutorToOrder successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SetActiveToOrder(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] SetActiveToOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orderID := ctx.Params("id")

	var IsActive models.ChangeActive

	if err := ctx.BodyParser(&IsActive); err != nil {
		custom_logger.Error("[OrderController] SetActiveToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := models.Valid.Struct(IsActive); err != nil {
		custom_logger.Error("[OrderController] SetActiveToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = c.OrderService.SetActiveOrderStatus(orderID, IsActive.IsActive, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] SetActiveToOrder successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SuggestOrder(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] SuggestOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	tutorID := ctx.Params("id")
	orderID := ctx.Query("order_id")

	err = c.OrderService.SuggestOrderToTutor(orderID, tutorID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] SuggestOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (c *OrderController) GetAllOrders(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] GetAllOrders called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orders, err := c.OrderService.GetAllOrders(UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] GetAllOrders successful")

	return ctx.JSON(orders)
}

func (c *OrderController) GetAllUsersOrders(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] GetAllUsersOrders called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orders, err := c.OrderService.GetAllUsersOrders(UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] GetAllUsersOrders successful")

	return ctx.JSON(orders)
}

func (c *OrderController) SetBanOrder(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] SetBanOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orderID := ctx.Params("id")

	err = c.OrderService.SetBanOrderStatus(orderID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] SetBanOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (c *OrderController) SetApprovedOrder(ctx *fiber.Ctx) error {
	custom_logger.Info("[OrderController] SetApprovedOrder called")

	UserData, err := getUserData(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}

	orderID := ctx.Params("id")

	err = c.OrderService.SetApprovedOrderStatus(orderID, UserData)

	if err != nil {
		code, currentError := utils.MapError(err)
		return ctx.Status(code).JSON(fiber.Map{"error": currentError})
	}

	custom_logger.Info("[OrderController] SetApprovedOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
