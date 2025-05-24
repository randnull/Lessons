package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	DBConfig
	ServerConfig
	BotConfig
	MQConfig
	GRPCConfig
}

type DBConfig struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
}

type ServerConfig struct {
	ServerPort string `env:"SERVER_PORT"`
	CorsOrigin string `env:"CORS_ORIGIN"`
}

type BotConfig struct {
	JWTSecret     string        `env:"JWT_SECRET"`
	BotToken      string        `env:"BOT_STUDENT_TOKEN"`
	BotTutorToken string        `env:"BOT_TUTOR_TOKEN"`
	AliveTime     time.Duration `env:"ALIVE_TIME"`
}

type MQConfig struct {
	User string `env:"MQ_USER"`
	Pass string `env:"MQ_PASS"`
	Host string `env:"MQ_HOST"`
	Port string `env:"MQ_PORT"`
}

type GRPCConfig struct {
	Host string `env:"GRPCUSERHOST"`
	Port string `env:"GRPCUSERPORT"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
