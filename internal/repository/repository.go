package repository

import (
	"github.com/pewpowder/url-shortener/pkg/logger"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db,
	}
}

func ConnectDatabase(DSN string, zl *zerolog.Logger, loggerCfg *logger.LoggerConfig) (*gorm.DB, error) {
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
