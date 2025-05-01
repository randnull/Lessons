package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig
	ServerConfig
	MQConfig
	SchedulerConfig
}

type DBConfig struct {
	DBHost     string `env:"DB_HOST"  env-default:"localhost"`
	DBPort     string `env:"DB_PORT" env-default:"5433"`
	DBUser     string `env:"DB_USER" env-default:"postgres"`
	DBPassword string `env:"DB_PASSWORD" env-default:"postgres"`
	DBName     string `env:"DB_NAME" env-default:"users_database"`
}

type MQConfig struct {
	User string `env:"MQ_USER" env-default:"guest"`
	Pass string `env:"MQ_PASS" env-default:"guest"`
	Host string `env:"MQ_HOST" env-default:"127.0.0.1"`
	Port string `env:"MQ_PORT" env-default:"5672"`
}

type ServerConfig struct {
	ServerPort string `env:"SERVER_PORT" env-default:"2000"`
}

type SchedulerConfig struct {
	Delay string `env:"SCHEDULER_RESPONSES_DELAY" env-default:"100"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
