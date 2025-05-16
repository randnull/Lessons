package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/pkg/custom_errors"
)

func ParseJWTToken(tokenStr string, jwtSecret string) (*models.Claims, error) {
	secretKey := []byte(jwtSecret)

	claims := &models.Claims{}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())

	token, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
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
