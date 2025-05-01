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
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
}

type MQConfig struct {
	User string `env:"MQ_USER"`
	Pass string `env:"MQ_PASS"`
	Host string `env:"MQ_HOST"`
	Port string `env:"MQ_PORT"`
}

type ServerConfig struct {
	ServerPort string `env:"SERVER_PORT"`
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
