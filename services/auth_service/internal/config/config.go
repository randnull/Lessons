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
	JWTsecret       string `env:"JWT_SECRET" env-default:"secret"`
	BotTokenStudent string `env:"BOT_TOKEN_STUDENT" env-default:"7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"`
	BotTokenTutor   string `env:"BOT_TOKEN_TUTOR" env-default:"7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"`
}
type GRPCConfig struct {
	Host string `env:"GRPCUSERHOST" env-default:"127.0.0.1"`
	Port string `env:"GRPCUSERPORT" env-default:"2000"`
}

type ServerConfig struct {
	ServerPort string `env:"SERVER_PORT" env-default:"8050"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
