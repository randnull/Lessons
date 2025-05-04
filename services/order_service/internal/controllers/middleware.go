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

		if UserClaims.Role != models.StudentType && UserClaims.Role != models.TutorType && UserClaims.Role != models.AdminType {
			logger.Debug("TokenAuthMiddleware failed. Get " + UserClaims.Role + " as Role")

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role forbidden",
			})
		}

		if userType == models.StudentType && UserClaims.Role != models.StudentType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role mismatch",
			})
		} else if userType == models.TutorType && UserClaims.Role != models.TutorType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role mismatch",
			})
		} else if userType == models.AdminType && UserClaims.Role != models.AdminType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Bad auth data provided. Role mismatch",
			})
		}

		if err != nil {
			logger.Info("TokenAuthMiddleware failed: " + err.Error())

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
