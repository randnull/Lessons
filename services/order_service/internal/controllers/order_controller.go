package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	initdata "github.com/telegram-mini-apps/init-data-golang"
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

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orderID, err := c.OrderService.CreateOrder(&order, InitData) // userData
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
		"orderID": orderID,
	})
}

// Get Info of current order by order_id
func (c *OrderController) GetOrderByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	order, err := c.OrderService.GetOrderById(id, InitData)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return ctx.JSON(order)
}

// Get all exist orders
func (c *OrderController) GetAllOrders(ctx *fiber.Ctx) error {
	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	//log.Printf("Controller: GetAllOrders")

	orders, err := c.OrderService.GetAllOrders(InitData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

// Get all order by user_id
func (c *OrderController) GetAllUsersOrders(ctx *fiber.Ctx) error {
	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	orders, err := c.OrderService.GetAllUsersOrders(InitData)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

// Delete order by id
func (c *OrderController) DeleteOrderByID(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	err := c.OrderService.DeleteOrder(orderID, InitData)

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

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "bad init data"})
	}

	err := c.OrderService.UpdateOrder(orderID, &order, InitData)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}
