package controllers

import (
	"github.com/gofiber/fiber/v2"
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
	var order models.NewOrder

	if err := ctx.BodyParser(&order); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.OrderService.CreateOrder(&order); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Order created successfully"})
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
