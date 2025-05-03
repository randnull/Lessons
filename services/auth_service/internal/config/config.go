package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	JWTConfig
	GRPCConfig
	ServerConfig
}

type JWTConfig struct {
	JWTsecret         string `env:"JWT_SECRET"`
	BotTokenStudent   string `env:"BOT_STUDENT_TOKEN"`
	BotTokenTutor     string `env:"BOT_TUTOR_TOKEN"`
	BotTokenAdmin     string `env:"BOT_ADMIN_TOKEN"`
	InitDataAliveTime int    `env:"INITDATA_ALIVE_TIME"`
	TokenAliveTime    int    `env:"TOKEN_ALIVE_TIME"`
	AdminId           int64  `env:"ADMIN_USER"`
}
type GRPCConfig struct {
	Host string `env:"GRPCUSERHOST"`
	Port string `env:"GRPCUSERPORT"`
}

type ServerConfig struct {
	ServerPort string `env:"AUTH_SERVICE_PORT"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
