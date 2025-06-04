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
	"github.com/pewpowder/url-shortener/internal/repository"
	"github.com/pewpowder/url-shortener/internal/resources"
	"github.com/pewpowder/url-shortener/internal/router"
	"github.com/pewpowder/url-shortener/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, env := config.MustLoad(resources.AppConfigFS)

	err := logger.InitLogger(resources.ZerologConfigFS, env)

	if err != nil {
		log.Fatalf("can't init logger: %s", err)
	}

	db, err := repository.ConnectDatabase(cfg.DB.DSN, logger.Get(), logger.GetConfig())

	if err != nil {
		logger.Get().Fatal().Err(err).Msg("failed to connect to database")
	}

	repo := repository.NewRepository(db)

	container := container.NewContainer(repo, cfg, env)

	g := gin.Default()

	router.Init(g, container)

	go g.Run(fmt.Sprintf(":%d", cfg.Server.Port))

	<-ctx.Done()
}
