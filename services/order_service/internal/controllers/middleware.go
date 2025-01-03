package controllers

import (
	"github.com/gofiber/fiber/v2"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"log"

	"github.com/randnull/Lessons/internal/config"
)

func TokenAuthMiddleware(cfg config.BotConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("token")
		log.Printf("token: %s", token)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "No token provided",
			})
		}

		userData, err := initdata.Parse(token)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad init data provided. Error in Parse",
			})
		}

		err = initdata.Validate(token, cfg.BotToken, cfg.AliveTime)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad init data provided. Error in Validate",
			})
		}

		c.Locals("user_data", userData)

		return c.Next()
	}
}
