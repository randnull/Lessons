package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerPort string `env:"PORT" env-default:"6050"`

	JWTsecret string `env:"JWT_SECRET" env-default:"secret"`
	BotToken  string `env:"BOT_TOKEN" env-default:"7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
