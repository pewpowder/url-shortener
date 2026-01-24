package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/container"
	"github.com/pewpowder/url-shortener/internal/resources"
	"github.com/pewpowder/url-shortener/pkg/logger"
	"github.com/rs/zerolog"
)

// TODO: remove gorm and gin. Use instead sqlx and http packages

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, env := config.MustLoad(resources.AppConfigFS)

	err := logger.InitLogger(resources.ZerologConfigFS, env)
	if err != nil {
		log.Fatalf("can't init logger: %s", err)
	}

	db, err := connectDatabase(ctx, cfg.DB.DSN, logger.Get(), logger.GetConfig())
	if err != nil {
		logger.Get().Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	container := container.NewContainer(db, cfg, env)
	_ = container

	// // gin.DebugPrintRouteFunc = logger.GinDebugPrintRoute
	// // gin.DebugPrintFunc = logger.GinDebugPrint
	// // g := gin.Default()

	// router.Init(g, container)

	// go g.Run(fmt.Sprintf(":%d", cfg.Server.Port))

	<-ctx.Done()
	
	// TODO: graceful shutdown
}

func connectDatabase(ctx context.Context, DSN string, zl *zerolog.Logger, loggerConfig *logger.LoggerConfig) (*sqlx.DB, error) {
	c, cancel := context.WithTimeout(ctx, time.Second * 10)
	defer cancel();
	
	cfg, err := pgx.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}
	
	traceLogger := logger.NewTraceLogger(*zl, loggerConfig.Pgx.LogLevel)
	m := logger.MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			otelpgx.NewTracer(),
			traceLogger,
		},
	}
	
	cfg.Tracer = &m;
	
	sqlDB := stdlib.OpenDB(*cfg)
	db := sqlx.NewDb(sqlDB, "pgx")
	
	if err := db.PingContext(c); err != nil {
		return nil, err
	}

	return db, nil
}

// TODO: Set up nginx for application
