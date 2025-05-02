package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/Lessons/internal/models"
	"time"
)

func CreateJWTToken(userID string, telegramID int64, username string, role string, jwtSecret string, TokenAlive int) (string, error) {
	secretKey := []byte(jwtSecret)

	claims := models.Claims{
		UserID:     userID,
		Username:   username,
		TelegramID: telegramID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(TokenAlive) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenStr, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
