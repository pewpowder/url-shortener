package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/resources"
	"github.com/pewpowder/url-shortener/pkg/logger"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, env := config.MustLoad(resources.AppConfigFS)

	err := logger.InitLogger(resources.ZerologConfigFS, env)

	if err != nil {
		log.Fatalf("can't init logger: %s", err)
	}

	connectDatabase(cfg.DB.DSN, logger.Get(), logger.GetConfig())

	if err != nil {
		logger.Get().Fatal().Err(err).Msg("failed to connect to database")
	}

	logger.Get().Info().Int("port", cfg.Server.Port).Msg("server started")

	<-ctx.Done()

	// if err := server.Shutdown(context.Background()); err != nil {
	// 	log.Fatalf("could not shutdown: %v\n", err)
	// }

	// storage postgreSQL
	// router chi, "'chi render' google it"
}

func connectDatabase(DSN string, zl *zerolog.Logger, loggerCfg *logger.LoggerConfig) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.NewGormLogger(zl, loggerCfg.ToGormConfig(&loggerCfg.Gorm)),
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), gormCfg) // TODO: Compare with default gorm logger in the future

	if err != nil {
		return nil, err
	}

	return db, nil
}
