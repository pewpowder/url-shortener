package main

import (
	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/resources"
	"github.com/pewpowder/url-shortener/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, env := config.MustLoad(resources.AppConfigFS)

	logger.InitLogger(resources.ZerologConfigFS, env)

	loggerCfg := logger.GetLoggerConfig(resources.ZerologConfigFS, env)
	gormCfg := &gorm.Config{
		Logger: logger.NewGormLogger(logger.Get(), loggerCfg.ToGormConfig(&loggerCfg.Gorm)),
	}

	_, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DB.DSN,
	}), gormCfg) // TODO: Compare with default gorm logger in the future

	if err != nil {
		logger.Get().Fatal().Msgf("Can't connect to DB: %s", err)
	}

	// storage postgreSQL
	// router chi, "'chi render' google it"
	// set graceful shutdown
}
