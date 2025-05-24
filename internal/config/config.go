package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"prod"`
	DB  DBConfig
}

type DBConfig struct {
	DSN string `env:"DSN" env-required:"true"`
}

type HTTPServer struct {
	Port        string        `env:"address" env-default:"8080"`
	Host        string        `env:"address" env-default:"localhost"`
	Timeout     time.Duration `env:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `env:"idle_timeout" env-default:"60s"`
}

func MustLoad(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Can't load config file: %s", err)
	}

	return &cfg
}
