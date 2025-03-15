package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/config"
	auth "github.com/randnull/Lessons/internal/jwt"
	"github.com/randnull/Lessons/internal/models"
)

func TokenAuthMiddleware(cfg config.BotConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "No token provided",
			})
		}
		fmt.Println(token, cfg.JWTSecret)
		UserClaims, err := auth.ParseJWTToken(token, cfg.JWTSecret)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Error in Parse",
			})
		}

		c.Locals("user_data",
			models.UserData{
				TelegramID: UserClaims.TelegramID,
				UserID:     UserClaims.UserID,
				Role:       UserClaims.Role,
			})

		return c.Next()
	}
}
