package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/Lessons/internal/custom_errors"
	"github.com/randnull/Lessons/internal/models"
)

func ParseJWTToken(tokenStr string, jwtSecret string) (*models.Claims, error) {
	secretKey := []byte(jwtSecret)

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !token.Valid {
		return nil, custom_errors.ErrorInvalidToken
	}

	return claims, nil
}
