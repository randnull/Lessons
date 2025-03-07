package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig
	ServerConfig
}

type DBConfig struct {
	DBHost     string `env:"DB_HOST"  env-default:"localhost"`     //  env-default:"dpg-cttubetumphs73eikdbg-a.oregon-postgres.render.com"
	DBPort     string `env:"DB_PORT" env-default:"5433"`           //  env-default:"5432"
	DBUser     string `env:"DB_USER" env-default:"postgres"`       //  env-default:"change"
	DBPassword string `env:"DB_PASSWORD" env-default:"postgres"`   // env-default:"9yuVZktnLKzqMrkywVgTlhDxVQsqWXbP"
	DBName     string `env:"DB_NAME" env-default:"users_database"` //  env-default:"orders_database_bhw2"
}

type ServerConfig struct {
	ServerPort string `env:"SERVER_PORT" env-default:"2000"` //  env-default:"6050"
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
