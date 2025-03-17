package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/models"
	"log"
	"time"
)

func CreateJWTToken(userID string, telegramID int64, username string, role models.RoleType, jwtSecret string) (string, error) {
	secretKey := []byte(jwtSecret)
	//
	//claims := jwt.MapClaims{
	//	"user_id":     userID,
	//	"telegram_id": telegramID,
	//	"role":        role,
	//}
	//claims := jwt.MapClaims{}

	claims := models.Claims{
		UserID:     userID,
		Username:   username,
		TelegramID: telegramID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30000 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenStr, err := token.SignedString(secretKey)

	if err != nil {
		log.Fatal(err)
	}

	return tokenStr, nil
}

func ParseJWTToken(tokenStr string, jwtSecret string) (*models.Claims, error) {
	secretKey := []byte(jwtSecret)

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, custom_errors.ErrorInvalidToken
	}

	return claims, nil
}
