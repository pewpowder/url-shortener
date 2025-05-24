package logger

import (
	"embed"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	gormLogger "gorm.io/gorm/logger"
)

type LoggerConfig struct {
	Zerolog    ZerologConfig `yaml:"zerolog"`
	Lumberjack Lumberjack    `yaml:"lumberjack"`
	Gorm       GormConfig    `yaml:"gorm"`
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

type GormConfig struct {
	SlowThreshold             time.Duration       `yaml:"slow_threshold"`
	Colorful                  bool                `yaml:"colorful"`
	IgnoreRecordNotFoundError bool                `yaml:"ignore_record_not_found_error"`
	ParameterizedQueries      bool                `yaml:"parameterized_queries"`
	LogLevel                  gormLogger.LogLevel `yaml:"log_level"`
}

var (
	once sync.Once
	zlog *zerolog.Logger
)

func InitLogger(yamlFS embed.FS, env string) {
	once.Do(func() {
		cfg := GetLoggerConfig(yamlFS, env)
		set(configureZerolog(cfg))
	})
}

func GetLoggerConfig(yamlFS embed.FS, env string) *LoggerConfig {
	file, err := yamlFS.ReadFile(fmt.Sprintf(config.ZEROLOG_CONFIG_PATH, env))

	if err != nil {
		log.Fatalf("Error during reading the logger config file: %s", err)
	}

	var cfg LoggerConfig

	err = yaml.Unmarshal(file, &cfg)

	if err != nil {
		log.Fatalf("Error during decoding the logger config file: %s", err)
	}

	return &cfg
}

func set(zerolog *zerolog.Logger) {
	if zlog != nil {
		return
	}

	zlog = zerolog
}

func Get() *zerolog.Logger {
	return zlog
}

func (l *LoggerConfig) ToGormConfig(gormCfg *GormConfig) *gormLogger.Config {
	return &gormLogger.Config{
		SlowThreshold:             gormCfg.SlowThreshold,
		Colorful:                  gormCfg.Colorful,
		IgnoreRecordNotFoundError: gormCfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      gormCfg.ParameterizedQueries,
		LogLevel:                  gormCfg.LogLevel,
	}
}
