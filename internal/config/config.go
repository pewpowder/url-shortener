package config

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port        int           `yaml:"port"`
	Host        string        `yaml:"host"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DBConfig struct {
	DSN string `yaml:"DSN"`
}

func (cfg *Config) validate() error {
	// mandatory fields
	if cfg.DB.DSN == "" {
		return errors.New("DSN must be setted")
	}

	if cfg.Server.Host == "" {
		return errors.New("server host must be setted")
	}

	if cfg.Server.Port == 0 {
		return errors.New("server port must be setted")
	}

	// default values
	if cfg.Server.Timeout == 0 {
		cfg.Server.Timeout = 5 * time.Second
	}

	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 5 * time.Second
	}

	return nil
}

func MustLoad(yamlFS embed.FS) (*Config, string) {
	var env *string
	if value := os.Getenv("env"); value != "" {
		env = &value
	} else {
		env = flag.String("env", DEV, "application enviroment")
		flag.Parse()
	}

	fullConfigPath := fmt.Sprintf(APP_CONFIG_PATH, *env)

	file, err := yamlFS.ReadFile(fullConfigPath)
	if err != nil {
		log.Fatalf("Can't read config file: %s", err)
	}

	var config Config

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Can't parse config file %s", err)
	}

	if err := config.validate(); err != nil {
		log.Fatalf("Invalid config file %s", err)
	}

	return &config, *env
}
