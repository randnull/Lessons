package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"

	//"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/service"
	"log"
)

type OrderController struct {
	OrderService service.OrderServiceInt
}

func NewOrderController(OrderService service.OrderServiceInt) *OrderController {
	return &OrderController{
		OrderService: OrderService,
	}
}

func (c *OrderController) HealtzHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// Create order
func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	var order models.NewOrder

	log.Printf("Controller: CreateOrder\nRequest model: %v", order)

	if err := ctx.BodyParser(&order); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("Create Order %+v", order)

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orderID, err := c.OrderService.CreateOrder(&order, UserData) // userData
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"orderID": orderID,
	})
}

func (c *OrderController) GetOrderByIdTutor(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	order, err := c.OrderService.GetOrderByIdTutor(id, UserData)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return ctx.JSON(order)
}

// Get Info of current order by order_id
func (c *OrderController) GetOrderByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	order, err := c.OrderService.GetOrderById(id, UserData)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return ctx.JSON(order)
}

// Get all exist orders
func (c *OrderController) GetAllOrders(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	//log.Printf("Controller: GetAllOrders")

	orders, err := c.OrderService.GetAllOrders(UserData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

func (c *OrderController) GetOrdersPagination(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})

	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})

	}

	tag := ctx.Query("tag")

	fmt.Println(tag)

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orders, err := c.OrderService.GetOrdersWithPagination(page, size, tag, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

func (c *OrderController) GetStudentOrdersPagination(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Page param not correct"})

	}

	size, err := strconv.Atoi(ctx.Query("size"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param not correct"})
	}

	if size > 100 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "size param more than maximum"})
	}

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orders, err := c.OrderService.GetStudentOrdersWithPagination(page, size, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

// Get all order by user_id
func (c *OrderController) GetAllUsersOrders(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	fmt.Println(UserData, ok)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orders, err := c.OrderService.GetAllUsersOrders(UserData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

// Delete order by id
func (c *OrderController) DeleteOrderByID(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	err := c.OrderService.DeleteOrder(orderID, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

// Update order by id .
func (c *OrderController) UpdateOrderByID(ctx *fiber.Ctx) error {

	log.Println("Пришел запрос на обновление:")
	orderID := ctx.Params("id")

	log.Printf("Update Order %+v", orderID)

	var order models.UpdateOrder

	fmt.Println(ctx)

	if err := ctx.BodyParser(&order); err != nil {
		log.Println("Error :", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("Update Order %+v", order)

	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	err := c.OrderService.UpdateOrder(orderID, &order, UserData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SelectTutorToOrder(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	responseID := ctx.Params("id")

	err := c.OrderService.SelectTutor(responseID, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get orders" + err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SetActiveToOrder(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orderID := ctx.Params("id")

	var IsActive models.ChangeActive

	if err := ctx.BodyParser(&IsActive); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request."})
	}

	err := c.OrderService.SetActiveOrderStatus(orderID, IsActive.IsActive, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to set inactive"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (c *OrderController) SuggestOrder(ctx *fiber.Ctx) error {
	UserData, ok := ctx.Locals("user_data").(models.UserData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	tutorID := ctx.Params("id")
	orderID := ctx.Query("order_id")

	err := c.OrderService.SuggestOrderToTutor(orderID, tutorID, UserData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to suggest order"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
