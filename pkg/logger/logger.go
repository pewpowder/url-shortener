package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

type LogFormat string

const (
	FormatJson LogFormat = "json"
	FormatText LogFormat = "text"
)

type LoggerConfig struct {
	Level  slog.Level
	Format LogFormat
	Output string
}

var Log *slog.Logger

func InitLogger(cfg *LoggerConfig) {
	var w io.Writer = os.Stdout
	if cfg.Output != "stdout" && cfg.Output != "" {
		f, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Can't open the file %s used fallback logger output", cfg.Output)
		}
		w = f
	}

	switch cfg.Format {
	case FormatJson:
		Log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: cfg.Level}))
	case FormatText:
		Log = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: cfg.Level}))
	default:
		Log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: cfg.Level}))
	}
}
