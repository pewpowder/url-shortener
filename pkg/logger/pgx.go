package logger

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type PgxConfig struct {
	LogLevel tracelog.LogLevel `yaml:"log_level"`
}

type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

type pgxLogger struct {
	logger zerolog.Logger
}

func (l pgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	logEvent := l.logger.With().Fields(data).Logger()
	zerologLevel := l.logger.GetLevel()
	if zerologLevel == zerolog.TraceLevel {
			level = tracelog.LogLevelTrace
	}

	switch level {
	case tracelog.LogLevelTrace:
			logEvent.Trace().Msg(msg)
	case tracelog.LogLevelDebug:
			logEvent.Debug().Msg(msg)
	case tracelog.LogLevelInfo:
			logEvent.Info().Msg(msg)
	case tracelog.LogLevelWarn:
			logEvent.Warn().Msg(msg)
	case tracelog.LogLevelError:
			logEvent.Error().Msg(msg)
	default:
			logEvent.Info().Msg(msg)
	}
}

func NewTraceLogger(zeroLog zerolog.Logger, level tracelog.LogLevel) *tracelog.TraceLog {
	return &tracelog.TraceLog{
			Logger:   pgxLogger{logger: zeroLog},
			LogLevel: level,
			Config: &tracelog.TraceLogConfig{
					TimeKey: "duration",
			},
	}
}

func (m *MultiQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, t := range m.Tracers {
			ctx = t.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

func (m *MultiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, t := range m.Tracers {
			t.TraceQueryEnd(ctx, conn, data)
	}
}