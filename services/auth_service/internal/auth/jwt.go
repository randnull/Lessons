package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/randnull/Lessons/internal/errors"
	"log"
)

func CreateJWTToken(user_id int64, jwt_secret string) (string, error) {
	secret_key := []byte(jwt_secret)

	claims := jwt.MapClaims{
		"user_id": user_id,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	token_str, err := token.SignedString(secret_key)

	if err != nil {
		log.Fatal(err)
	}

	return token_str, nil
}

func ParseJWTToken(token_str string, jwt_secret string) (string, error) {
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.InvalidToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	str := fmt.Sprintf("%v", claims["user_id"])

	user_id := str

	return user_id, nil
}
