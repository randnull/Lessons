// Сервис для аутификации
// Цель - получить initdata, проверить ее на корректность
// Если все ок - распарсить ее и передать в mock сервиса UserService (grpc)
// Вернуть jwt токен с информации о пользователе (он содержит user_id)

// При запросе в другой сервис фронт передает jwt токен в него, там он проверяется на корректность
// Если все ок - запрос происходит. внутри токена - user_id и роль (репет, ученик)

// Конфиг передаватель в auth service

// Ручки endpoint просто принимают один запрос - login. в случае, если акк уже есть - то просто вернется jwt. если акк еще нет - он создаться и вернется jwt

package service

import (
	"fmt"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"time"

	"github.com/randnull/Lessons/internal/auth"
	"github.com/randnull/Lessons/internal/config"
)

type AuthServiceInt interface {
	Login(initdata string) (string, error)
}

type AuthService struct {
	cfg *config.Config
}

func NewAuthService(cfg *config.Config) AuthServiceInt {
	return &AuthService{
		cfg: cfg,
	}
}

func (authserv *AuthService) Login(initdata_from_user string) (string, error) {
	user_data, err := initdata.Parse(initdata_from_user)

	fmt.Println(user_data)

	if err != nil {
		return "", err
	}

	// REQ TO USER SERV user_id

	fmt.Println(initdata.Validate(initdata_from_user, authserv.cfg.BotToken, time.Hour))

	err = initdata.Validate(initdata_from_user, authserv.cfg.BotToken, time.Hour) // конфиг

	if err != nil {
		return "", err
	}

	jwt_token, err := auth.CreateJWTToken(user_data.User.ID, authserv.cfg.JWTsecret)
	if err != nil {
		return "", err
	}

	fmt.Println(jwt_token)

	return jwt_token, nil
}
