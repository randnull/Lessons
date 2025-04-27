package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/internal/config"
	auth "github.com/randnull/Lessons/internal/jwt"
	"github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
)

func TokenAuthMiddleware(cfg config.BotConfig, userType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		logger.Debug("TokenAuthMiddleware called")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "No token provided",
			})
		}

		UserClaims, err := auth.ParseJWTToken(token, cfg.JWTSecret)

		if err != nil {
			logger.Debug("TokenAuthMiddleware failed: " + err.Error())

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Error in Parse",
			})
		}

		if UserClaims.Role != "Student" && UserClaims.Role != "Tutor" && UserClaims.Role != "Admin" {
			logger.Debug("TokenAuthMiddleware failed. Get " + UserClaims.Role + " as Role")

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role forbidden",
			})
		}

		if userType == "Student" && UserClaims.Role != "Student" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role mismatch",
			})
		} else if userType == "Tutor" && UserClaims.Role != "Tutor" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role mismatch",
			})
		} else if userType == "Admin" && UserClaims.Role != "Admin" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role mismatch",
			})
		}

		if err != nil {
			logger.Debug("TokenAuthMiddleware failed: " + err.Error())

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Error in Parse",
			})
		}

		c.Locals("user_data",
			models.UserData{
				TelegramID: UserClaims.TelegramID,
				Username:   UserClaims.Username,
				UserID:     UserClaims.UserID,
				Role:       UserClaims.Role,
			})
		logger.Debug("TokenAuthMiddleware Success")

		return c.Next()
	}
}
