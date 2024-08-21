package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Env string `env:"ENV" env-default:"local"`
	HTTPServer
	Minio
}

type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS"  env-default:"localhost:8080"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT"  env-default:"4s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type Minio struct {
	Endpoint  string `env:"MINIO_ENDPOINT" env-default:"minio:9000"`
	AccessKey string `env:"MINIO_ROOT_USER" env-default:"minioadmin"`
	Password  string `env:"MINIO_ROOT_PASSWORD" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
