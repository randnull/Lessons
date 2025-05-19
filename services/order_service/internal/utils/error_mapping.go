package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/randnull/Lessons/pkg/custom_errors"
)

func MapError(err error) (int, error) {
	switch {
	case errors.Is(err, custom_errors.ErrorServiceError):
		return fiber.StatusInternalServerError, err
	case errors.Is(err, custom_errors.ErrorNotFound):
		return fiber.StatusNotFound, err
	case errors.Is(err, custom_errors.ErrNotAllowed):
		return fiber.StatusForbidden, err
	case errors.Is(err, custom_errors.ErrorBadStatus):
		return fiber.StatusBadRequest, err
	case errors.Is(err, custom_errors.ErrorParams):
		return fiber.StatusBadRequest, err
	case errors.Is(err, custom_errors.ErrorBanWords):
		return fiber.StatusBadRequest, err
	}

	return fiber.StatusInternalServerError, errors.New("unknown error happens")
}
