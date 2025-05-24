package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	traceErrMsg  = "%s %s\n[%.3fms] [rows:%v] %s"
	traceWarnMsg = "%s %s\n[%.3fms] [rows:%v] %s"
	traceInfoMsg = "%s\n[%.3fms] [rows:%v] %s"
)

type GormLogger struct {
	zerolog *zerolog.Logger
	config  *gormLogger.Config
}

func NewGormLogger(zl *zerolog.Logger, cfg *gormLogger.Config) gormLogger.Interface {
	logger := zl.With().Str("from", "gorm").Logger()

	return &GormLogger{
		zerolog: &logger,
		config:  cfg,
	}
}

func (l GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l GormLogger) Info(_ context.Context, msg string, opts ...any) {
	if l.config.LogLevel >= gormLogger.Info {
		l.zerolog.Info().Msgf(msg, opts...)
	}
}

func (l GormLogger) Warn(_ context.Context, msg string, opts ...any) {
	if l.config.LogLevel >= gormLogger.Warn {
		l.zerolog.Warn().Msgf(msg, opts...)
	}
}

func (l GormLogger) Error(_ context.Context, msg string, opts ...any) {
	if l.config.LogLevel >= gormLogger.Error {
		l.zerolog.Error().Msgf(msg, opts...)
	}
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	if l.config.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	slowThreshold := l.config.SlowThreshold * time.Microsecond
	sql, rows := f()

	/*
		Проблема: В логах не хватает информации о контексте запроса (например, request id, user id).
		Рекомендация: Прокидывать дополнительные поля из context.Context, если они есть.
	*/

	switch {
	case err != nil && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.config.IgnoreRecordNotFoundError):
		// TODO: check how will be zerolog log caller from gorm
		l.Error(ctx, traceErrMsg, utils.FileWithLineNum(), err, elapsed.Seconds()*1000, rows, sql)
	case elapsed >= slowThreshold:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold)
		l.Warn(ctx, traceWarnMsg, utils.FileWithLineNum(), slowLog, elapsed.Seconds()*1000, rows, sql)
	default:
		l.Info(ctx, traceInfoMsg, utils.FileWithLineNum(), elapsed.Seconds()*1000, rows, sql)
	}
}
