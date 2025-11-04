package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/resources"
	"github.com/pewpowder/url-shortener/internal/router"
	"github.com/pewpowder/url-shortener/pkg/logger"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	c, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, env := config.MustLoad(resources.AppConfigFS)

	err := logger.InitLogger(resources.ZerologConfigFS, env)

	if err != nil {
		log.Fatalf("can't init logger: %s", err)
	}

	db, err := connectDatabase(cfg.DB.DSN, logger.Get(), logger.GetConfig())

	if err != nil {
		logger.Get().Fatal().Err(err).Msg("failed to connect to database")
	}

	container := container.NewContainer(db, cfg, env)

	gin.DebugPrintRouteFunc = logger.GinDebugPrintRoute
	gin.DebugPrintFunc = logger.GinDebugPrint
	g := gin.Default()

	router.Init(g, container)

	go g.Run(fmt.Sprintf(":%d", cfg.Server.Port))

	<-c.Done()
}

func connectDatabase(DSN string, zl *zerolog.Logger, loggerCfg *logger.LoggerConfig) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger:         logger.NewGormLogger(zl, loggerCfg.ToGormConfig(&loggerCfg.Gorm)),
		TranslateError: true,
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), gormCfg) // TODO: Compare with default gorm logger in the future

	if err != nil {
		return nil, err
	}

	return db, nil
}

// TODO: Set up nginx for application!!!!
