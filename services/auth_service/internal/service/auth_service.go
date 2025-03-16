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
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/auth"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"log"
)

type AuthServiceInt interface {
	Login(AuthData *models.AuthData) (string, error)
}

type AuthService struct {
	cfg        *config.JWTConfig
	GRPCClient gRPC_client.GRPCClientInt
}

func NewAuthService(cfg *config.JWTConfig, grpcClient gRPC_client.GRPCClientInt) AuthServiceInt {
	return &AuthService{
		cfg:        cfg,
		GRPCClient: grpcClient,
	}
}

func (authserv *AuthService) Login(AuthData *models.AuthData) (string, error) {
	userData, err := initdata.Parse(AuthData.InitData)

	fmt.Println(userData)

	if err != nil {
		return "", err
	}

	// ВАЛИДАЦИЮ ОБЯЗАТЕЛЬНО ВКЛЮЧИТЬ КОГДА ОПРЕДЕЛИМСЯ С ТОКЕНАМИ
	//var errValidate error

	//switch AuthData.Role {
	//case models.RoleTutor:
	//	errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenTutor, time.Hour*30000) // конфиг
	//case models.RoleStudent:
	//	errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenStudent, time.Hour*30000) // конфиг
	//}

	//if errValidate != nil {
	//	return "", errValidate
	//}
	// create user создает или возвращает пользователя

	log.Println(AuthData.Role)

	//if userData.User.Username not exist -> error
	
	userID, err := authserv.GRPCClient.CreateUser(context.Background(), &models.NewUser{
		TelegramID: userData.User.ID,
		Username:   userData.User.Username,
		Name:       userData.User.FirstName,
		Role:       AuthData.Role,
	})

	if err != nil {
		return "", err
	}

	jwtToken, err := auth.CreateJWTToken(userID, userData.User.ID, AuthData.Role, authserv.cfg.JWTsecret)
	if err != nil {
		return "", err
	}
	x, _ := auth.ParseJWTToken(jwtToken, authserv.cfg.JWTsecret)
	log.Println("decoded", x)

	return jwtToken, nil
}
