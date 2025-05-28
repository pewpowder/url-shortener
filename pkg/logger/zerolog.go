package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"log"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func configureZerolog(cfg *LoggerConfig) (*zerolog.Logger, error) {
	writers, err := getWriters(cfg)

	if err != nil {
		return nil, err
	}

	zl := zerolog.New(io.MultiWriter(writers...))

	if cfg.Zerolog.TimeFormat != "" {
		zerolog.TimeFieldFormat = cfg.Zerolog.TimeFormat
	}

	if cfg.Zerolog.Level >= zerolog.TraceLevel && cfg.Zerolog.Level <= zerolog.Disabled {
		zerolog.SetGlobalLevel(cfg.Zerolog.Level)
	}

	if cfg.Zerolog.WithTimestamp {
		zl = zl.With().Timestamp().Logger()
	}

	if cfg.Zerolog.WithCaller {
		zl = zl.With().Caller().Logger()
	}

	return &zl, nil
}

func getWriters(cfg *LoggerConfig) ([]io.Writer, error) {
	var writers []io.Writer

	for _, output := range cfg.Zerolog.Outputs {
		switch output.Type {
		case "console":
			writers = append(writers, getConsoleWriter(&cfg.Zerolog, output.Name))
		case "file":
			writer, err := getFileWriter(&cfg.Lumberjack, output.Name)

			if err != nil {
				return nil, err
			}

			writers = append(writers, writer)
		default:
			log.Printf("Unknown output type: %s. Defaulting to console.", output)
			writers = append(writers, getConsoleWriter(&cfg.Zerolog, "stdout"))
		}
	}

	return writers, nil
}

func getConsoleWriter(cfg *ZerologConfig, name string) io.Writer {
	var out io.Writer

	if name == "stdout" {
		out = os.Stdout
	}

	if name == "stderr" {
		out = os.Stderr
	}

	if cfg.BeautifyLogs {
		return zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: cfg.TimeFormat,
		}
	}

	return out
}

func getFileWriter(rotateCfg *Lumberjack, name string) (io.Writer, error) {
	if err := os.MkdirAll(rotateCfg.Directory, 0755); err != nil {
		return nil, fmt.Errorf("can't create the directory (%s) for a log file: %w", rotateCfg.Directory, err)
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(rotateCfg.Directory, name),
		MaxSize:    rotateCfg.MaxSize,
		MaxAge:     rotateCfg.MaxAge,
		MaxBackups: rotateCfg.MaxBackups,
		Compress:   rotateCfg.Compress,
	}, nil
}
