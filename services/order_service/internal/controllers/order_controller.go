package controllers

import (
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

func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	var order models.NewOrder

	if err := ctx.BodyParser(&order); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Printf("Create Order %+v", order)

	InitData, ok := ctx.Locals("user_data").(initdata.InitData)

	if !ok { // убрать!
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
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

func (c *OrderController) GetOrderByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	order, err := c.OrderService.GetOrderById(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return ctx.JSON(order)
}

func (c *OrderController) GetAllOrders(ctx *fiber.Ctx) error {
	orders, err := c.OrderService.GetAllOrders()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get orders"})
	}

	return ctx.JSON(orders)
}

func (c *OrderController) DeleteOrderByID(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (c *OrderController) UpdateOrderByID(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
