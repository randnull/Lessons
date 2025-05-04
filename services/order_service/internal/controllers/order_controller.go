package controllers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/logger"
	"strconv"

	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
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
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Info("CreateOrder called. UserID: " + UserData.UserID + ", Role: " + UserData.Role)

	var order models.NewOrder

	if err := ctx.BodyParser(&order); err != nil {
		logger.Error("CreateOrder failed to parse body: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := models.Valid.Struct(order); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	orderID, err := c.OrderService.CreateOrder(&order, UserData)

	if err != nil {
		if errors.Is(err, custom_errors.ErrorServiceError) {
			logger.Error("CreateOrder failed: " + err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("CreateOrder successful, OrderID: " + orderID)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"orderID": orderID,
	})
}

func (c *OrderController) GetOrderByIdTutor(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	id := ctx.Params("id")

	logger.Info(fmt.Sprintf("GetOrderByIdTutor called. UserID: %v Role: %v OrderID: %v", UserData.UserID, UserData.Role, id))

	order, err := c.OrderService.GetOrderByIdTutor(id, UserData)

	if err != nil {
		if errors.Is(err, custom_errors.ErrorServiceError) {
			logger.Error("GetOrderByIdTutor failed: " + err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
		}
		if errors.Is(err, custom_errors.ErrorNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("GetOrderByIdTutor successful")

	return ctx.Status(fiber.StatusOK).JSON(order)
}

func (c *OrderController) GetOrderByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	logger.Info(fmt.Sprintf("GetOrderByID called. UserID: %v Role: %v OrderID: %v", UserData.UserID, UserData.Role, id))

	order, err := c.OrderService.GetOrderById(id, UserData)

	if err != nil {
		if errors.Is(err, custom_errors.ErrorServiceError) {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
		}
		if errors.Is(err, custom_errors.ErrNotAllowed) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not allowed"})
		}
		if errors.Is(err, custom_errors.ErrorNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("GetOrderByID successful")

	return ctx.Status(fiber.StatusOK).JSON(order)
}

func (c *OrderController) GetOrdersPagination(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page"))

	logger.Info("GetOrdersPagination called")

	if err != nil {
		logger.Error(fmt.Sprintf("GetOrdersPagination failed: %v", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})
	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		logger.Error(fmt.Sprintf("GetOrdersPagination failed: %v", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	tag := ctx.Query("tag")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	orders, err := c.OrderService.GetOrdersWithPagination(page, size, tag, UserData)

	if err != nil {
		logger.Error(fmt.Sprintf("GetOrdersPagination failed: %v", err.Error()))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
	}
	logger.Info("GetOrdersPagination successful")

	return ctx.Status(fiber.StatusOK).JSON(orders)
}

func (c *OrderController) GetStudentOrdersPagination(ctx *fiber.Ctx) error {
	logger.Info("GetStudentOrdersPagination called")

	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		logger.Error(fmt.Sprintf("GetStudentOrdersPagination failed: %v", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})
	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		logger.Error(fmt.Sprintf("GetStudentOrdersPagination failed: %v", err.Error()))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	if size > 100 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param more than maximum"})
	}

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	orders, err := c.OrderService.GetStudentOrdersWithPagination(page, size, UserData)

	if err != nil {
		logger.Error("GetStudentOrdersPagination failed: " + err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
	}

	logger.Info("GetStudentOrdersPagination successful")

	return ctx.Status(fiber.StatusOK).JSON(orders)
}

func (c *OrderController) DeleteOrderByID(ctx *fiber.Ctx) error {
	logger.Info("DeleteOrderByID called")

	orderID := ctx.Params("id")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	err := c.OrderService.DeleteOrder(orderID, UserData)

	if err != nil {
		if errors.Is(err, custom_errors.ErrorServiceError) {
			logger.Error(fmt.Sprintf("DeleteOrderByID failed: %v", err.Error()))
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
		}
		if errors.Is(err, custom_errors.ErrNotAllowed) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not allowed"})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Info("DeleteOrderByID successful")

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

// вот тут закончил

func (c *OrderController) UpdateOrderByID(ctx *fiber.Ctx) error {
	logger.Info("UpdateOrderByID called")

	orderID := ctx.Params("id")

	var order models.UpdateOrder

	if err := ctx.BodyParser(&order); err != nil {
		logger.Error("UpdateOrderByID failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	err := c.OrderService.UpdateOrder(orderID, &order, UserData)
	if err != nil {
		logger.Error("UpdateOrderByID failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	logger.Info("UpdateOrderByID successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SelectTutorToOrder(ctx *fiber.Ctx) error {
	logger.Info("SelectTutorToOrder called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	responseID := ctx.Params("id")

	err := c.OrderService.SelectTutor(responseID, UserData)

	if err != nil {
		logger.Error("SelectTutorToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get orders" + err.Error()})
	}

	logger.Info("SelectTutorToOrder successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SetActiveToOrder(ctx *fiber.Ctx) error {
	logger.Debug("SetActiveToOrder called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	orderID := ctx.Params("id")

	var IsActive models.ChangeActive

	if err := ctx.BodyParser(&IsActive); err != nil {
		logger.Error("SetActiveToOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	if err := models.Valid.Struct(IsActive); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := c.OrderService.SetActiveOrderStatus(orderID, IsActive.IsActive, UserData)

	if err != nil {
		logger.Error("SetActiveToOrder failed: " + err.Error())

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to set inactive"})
	}

	logger.Info("SetActiveToOrder successful")

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SuggestOrder(ctx *fiber.Ctx) error {
	logger.Debug("SuggestOrder called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	tutorID := ctx.Params("id")
	orderID := ctx.Query("order_id")

	err := c.OrderService.SuggestOrderToTutor(orderID, tutorID, UserData)

	if err != nil {
		logger.Error("SuggestOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to suggest order"})
	}

	logger.Info("SuggestOrder successful")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (c *OrderController) GetAllOrders(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	orders, err := c.OrderService.GetAllOrders(UserData)

	logger.Debug("GetAllOrders called. UserID: " + UserData.UserID + ", Role: " + UserData.Role)

	if err != nil {
		logger.Error("GetAllOrders failed: " + err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	logger.Info("GetAllOrders successful")

	return ctx.JSON(orders)
}

func (c *OrderController) GetAllUsersOrders(ctx *fiber.Ctx) error {
	logger.Debug("GetAllUsersOrders called")

	UserData, _ := ctx.Locals("user_data").(models.UserData)

	orders, err := c.OrderService.GetAllUsersOrders(UserData)

	if err != nil {
		logger.Error("GetAllUsersOrders failed: " + err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
	}

	logger.Info("GetAllUsersOrders successful")

	return ctx.JSON(orders)
}

func (c *OrderController) SetBanOrder(ctx *fiber.Ctx) error {
	UserData, _ := ctx.Locals("user_data").(models.UserData)

	orderID := ctx.Params("id")

	logger.Debug(fmt.Sprintf("SetBanOrder called. UserID: %v OrderID: %v", UserData.UserID, orderID))

	err := c.OrderService.SetBanOrderStatus(orderID, UserData)

	if err != nil {
		logger.Error("SetBanOrder failed: " + err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Service error"})
	}

	logger.Info(fmt.Sprintf("SetBanOrder successful for OrderID: %v", orderID))

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
