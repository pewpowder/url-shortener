package logger

import (
	"embed"
	"fmt"
	"sync"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type LoggerConfig struct {
	Zerolog    ZerologConfig `yaml:"zerolog"`
	Lumberjack Lumberjack    `yaml:"lumberjack"`
	Pgx PgxConfig `yaml:"pgx"`
}

type ZerologConfig struct {
	Level         zerolog.Level `yaml:"level"`
	Outputs       []LogOutput   `yaml:"outputs"`
	TimeFormat    string        `yaml:"time_format"`
	WithTimestamp bool          `yaml:"with_timestamp"`
	WithCaller    bool          `yaml:"with_caller"`
	BeautifyLogs  bool          `yaml:"beautify_logs"`
}

type LogOutput struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

type Lumberjack struct {
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
	Directory  string `yaml:"directory"`
}

var (
	once      sync.Once
	zlog      *zerolog.Logger
	loggerCfg *LoggerConfig
)

func InitLogger(yamlFS embed.FS, env string) error {
	var initError error
	once.Do(func() {
		cfg, err := getLoggerConfig(yamlFS, env)
		if err != nil {
			initError = err
			return
		}

		zl, err := configureZerolog(cfg)
		if err != nil {
			initError = err
			return
		}

		setLogger(zl)
		setLoggerConfig(cfg)
	})
	return initError
}

func getLoggerConfig(yamlFS embed.FS, env string) (*LoggerConfig, error) {
	file, err := yamlFS.ReadFile(fmt.Sprintf(config.ZEROLOG_CONFIG_PATH, env))
	if err != nil {
		return nil, fmt.Errorf("can't read logger config: %w", err)
	}

	var cfg LoggerConfig
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't decode logger config: %w", err)
	}

	return &cfg, nil
}

func setLogger(zerolog *zerolog.Logger) {
	if zlog != nil {
		return
	}

	zlog = zerolog
}

func setLoggerConfig(cfg *LoggerConfig) {
	if loggerCfg != nil {
		return
	}

	loggerCfg = cfg
}

func Get() *zerolog.Logger {
	return zlog
}

func GetConfig() *LoggerConfig {
	return loggerCfg
}
