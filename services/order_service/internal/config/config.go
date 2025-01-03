package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	DBConfig
	ServerConfig
	BotConfig
}

type DBConfig struct {
	DBHost     string `env:"DB_HOST" env-default:"localhost"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" env-default:"postgres"`
	DBName     string `env:"DB_NAME" env-default:"postgres"`
}

type ServerConfig struct {
	ServerPort string `env:"PORT" env-default:"6050"`
}

type BotConfig struct {
	BotToken  string        `env:"BOT_TOKEN" env-default:"7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"`
	AliveTime time.Duration `env:"ALIVE_TIME" env-default:"30h"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
