package service

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/auth"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/gRPC_client"
	lg "github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"time"
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
	lg.Info("parsing jwt")

	userData, err := initdata.Parse(AuthData.InitData)

	if err != nil {
		lg.Error(fmt.Sprintf("jwt parsing create error: %v", err))
		return "", err
	}

	lg.Info(fmt.Sprintf("parsing ok. User-Telegram-Id: %v. User-Role: %v", userData.User.ID, AuthData.Role))

	aliveTime := authserv.cfg.InitDataAliveTime

	var errValidate error

	switch AuthData.Role {
	case models.RoleTutor:
		errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenTutor, time.Duration(aliveTime)*time.Minute)
	case models.RoleStudent:
		errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenStudent, time.Duration(aliveTime)*time.Minute)
	}

	if errValidate != nil {
		lg.Error(fmt.Sprintf("Error validation. User-Telegram-Id: %v. User-Role: %v. Error: %v", userData.User.ID, AuthData.Role, errValidate.Error()))
		return "", errValidate
	}

	lg.Info("request to create user")

	userID, err := authserv.GRPCClient.CreateUser(context.Background(), &models.NewUser{
		TelegramID: userData.User.ID,
		Name:       userData.User.FirstName,
		Role:       AuthData.Role,
	})

	lg.Info(fmt.Sprintf("User Created ok. User-Telegram-Id: %v. User-Role: %v. User-Id: %v", userData.User.ID, AuthData.Role, userID))

	if err != nil {
		lg.Error(fmt.Sprintf("Error Create User. User-Telegram-Id: %v. User-Role: %v. Error: %v", userData.User.ID, AuthData.Role, err.Error()))
		return "", err
	}

	lg.Info(fmt.Sprintf("Trying create jwt token. User-Telegram-Id: %v. User-Role: %v. User-Id: %v", userData.User.ID, AuthData.Role, userID))
	jwtToken, err := auth.CreateJWTToken(userID, userData.User.ID, userData.User.Username, AuthData.Role, authserv.cfg.JWTsecret, authserv.cfg.TokenAliveTime)

	if err != nil {
		lg.Error(fmt.Sprintf("Error create jwt token. User-Telegram-Id: %v. User-Role: %v. User-Id: %v. Error: %v", userData.User.ID, AuthData.Role, userID, err.Error()))
		return "", err
	}

	lg.Info(fmt.Sprintf("jwt token ok. User-Telegram-Id: %v. User-Role: %v. User-Id: %v", userData.User.ID, AuthData.Role, userID))

	return jwtToken, nil
}
