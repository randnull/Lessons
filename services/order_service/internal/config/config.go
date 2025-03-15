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
	DBHost     string `env:"DB_HOST" env-default:"localhost"`       //  env-default:"dpg-cttubetumphs73eikdbg-a.oregon-postgres.render.com"
	DBPort     string `env:"DB_PORT" env-default:"5432"`            //  env-default:"5432"
	DBUser     string `env:"DB_USER" env-default:"postgres"`        //  env-default:"change"
	DBPassword string `env:"DB_PASSWORD" env-default:"postgres"`    // env-default:"9yuVZktnLKzqMrkywVgTlhDxVQsqWXbP"
	DBName     string `env:"DB_NAME" env-default:"orders_database"` //  env-default:"orders_database_bhw2"
}

type ServerConfig struct {
	ServerPort string `env:"SERVER_PORT" env-default:"6050"` //  env-default:"6050"
}

type BotConfig struct {
	JWTSecret string        `env:"JWT_TOKEN" env-default:"secret"`
	BotToken  string        `env:"BOT_TOKEN" env-default:"7629903300:AAFwHNldwaNDI8cqv7FneC6DtYetbhe0DP0"`
	AliveTime time.Duration `env:"ALIVE_TIME" env-default:"30h"`
}

type MQConfig struct {
	User string `env:"MQ_USER" env-default:"guest"`     //  env-default:"guest"
	Pass string `env:"MQ_PASS" env-default:"guest"`     //  env-default:"guest"
	Host string `env:"MQ_HOST" env-default:"127.0.0.1"` //  env-default:"rabbitmq"
	Port string `env:"MQ_PORT" env-default:"5672"`      // env-default:"5672"
}

type GRPCConfig struct {
	Host string `env:"GRPCUSERHOST" env-default:"127.0.0.1"`
	Port string `env:"GRPCUSERPORT" env-default:"2000"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
