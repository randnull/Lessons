package service

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/auth"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/gRPC_client"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/pkg/custom_errors"
	lg "github.com/randnull/Lessons/pkg/logger"
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
	lg.Info(fmt.Sprintf("Login started , role=%v", AuthData.Role))

	userData, err := initdata.Parse(AuthData.InitData)
	if err != nil {
		lg.Error(fmt.Sprintf("InitData parse failed, error=%v", err))
		return "", err
	}

	lg.Info(fmt.Sprintf("InitData parsed, telegram_id=%v, role=%v", userData.User.ID, AuthData.Role))

	aliveTime := time.Duration(authserv.cfg.InitDataAliveTime) * time.Hour

	var errValidate error

	switch AuthData.Role {
	case models.RoleTutor:
		errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenTutor, aliveTime)
	case models.RoleStudent:
		errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenStudent, aliveTime)
	case models.RoleAdmin:
		errValidate = initdata.Validate(AuthData.InitData, authserv.cfg.BotTokenAdmin, aliveTime)
	default:
		errValidate = custom_errors.ErrorInvalidRole
	}

	if errValidate != nil {
		lg.Error(fmt.Sprintf("Validation failed. telegram_id=%v, role=%v, error=%v", userData.User.ID, AuthData.Role, errValidate))
		return "", errValidate
	}

	if AuthData.Role == models.RoleAdmin && userData.User.ID != authserv.cfg.AdminId {
		lg.Info(fmt.Sprintf("Unauthorized admin login attempt, telegram_id=%v", userData.User.ID))
		return "", custom_errors.ErrorInvalidToken
	}

	if AuthData.Role != models.RoleAdmin {
		user, err := authserv.GRPCClient.GetUserByTelegramID(context.Background(), userData.User.ID, AuthData.Role)
		if err != nil {
			lg.Error(fmt.Sprintf("GetUserByTelegramID failed, telegram_id=%v, role=%v, error=%v", userData.User.ID, AuthData.Role, err))
		} else if user.IsBanned {
			lg.Info(fmt.Sprintf("Banned user login attempt, telegram_id=%v, role=%v", userData.User.ID, AuthData.Role))
			return "", custom_errors.ErrorInvalidToken
		}
	}

	lg.Info(fmt.Sprintf("Creating user, telegram_id=%v, role=%v", userData.User.ID, AuthData.Role))

	userID, err := authserv.GRPCClient.CreateUser(context.Background(), &models.NewUser{
		TelegramID: userData.User.ID,
		Name:       userData.User.FirstName,
		Role:       AuthData.Role,
	})

	if err != nil {
		lg.Error(fmt.Sprintf("CreateUser failed, telegram_id=%v, role=%v, error=%v", userData.User.ID, AuthData.Role, err))
		return "", err
	}

	lg.Info(fmt.Sprintf("User created, telegram_id=%v, role=%v, user_id=%v", userData.User.ID, AuthData.Role, userID))

	lg.Info(fmt.Sprintf("Creating JWT token, telegram_id=%v, role=%v, user_id=%v", userData.User.ID, AuthData.Role, userID))

	jwtToken, err := auth.CreateJWTToken(userID, userData.User.ID, userData.User.Username, AuthData.Role, authserv.cfg.JWTsecret, authserv.cfg.TokenAliveTime)

	if err != nil {
		lg.Error(fmt.Sprintf("JWT token creation failed, telegram_id=%v, user_id=%v, role=%v, error=%v", userData.User.ID, userID, AuthData.Role, err))
		return "", err
	}

	lg.Info(fmt.Sprintf("Login successful, telegram_id=%v, role=%v, user_id=%v", userData.User.ID, AuthData.Role, userID))

	return jwtToken, nil
}
