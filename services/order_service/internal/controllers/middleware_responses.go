package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/config"
	auth "github.com/randnull/Lessons/internal/jwt"
	"github.com/randnull/Lessons/internal/models"
)

func TokenAuthMiddlewareResponses(cfg config.BotConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "No token provided",
			})
		}

		UserClaims, err := auth.ParseJWTToken(token, cfg.JWTSecret)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Error in Parse",
			})
		}
		//
		//if UserClaims.Role != "Tutor" {
		//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		//		"error": true,
		//		"msg":   "Not tutor",
		//	})
		//}

		c.Locals("user_data",
			models.UserData{
				TelegramID: UserClaims.TelegramID,
				Username:   UserClaims.Username,
				UserID:     UserClaims.UserID,
				Role:       UserClaims.Role,
			})

		return c.Next()
	}
}
