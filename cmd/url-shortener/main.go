package main

import (
	"log"
	"log/slog"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad(".env")

	_, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DB.DSN,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Can't connect to DB: %s", err)
	}

	logger.InitLogger(&logger.LoggerConfig{
		Level:  slog.LevelInfo,
		Format: logger.FormatJson,
		Output: "stdout",
	})

	// storage postgreSQL
	// router chi, "'chi render' google it"
}
