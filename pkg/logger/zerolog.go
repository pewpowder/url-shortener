package logger

import (
	"io"
	"os"
	"path/filepath"

	"log"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func configureZerolog(cfg *LoggerConfig) *zerolog.Logger {
	writers := getWriters(cfg)
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

	return &zl
}

func getWriters(cfg *LoggerConfig) []io.Writer {
	var writers []io.Writer

	for _, output := range cfg.Zerolog.Outputs {
		switch output.Type {
		case "console":
			writers = append(writers, getConsoleWriter(&cfg.Zerolog, output.Name))
		case "file":
			writers = append(writers, getFileWriter(&cfg.Lumberjack, output.Name))
		default:
			log.Printf("Unknown output type: %s. Defaulting to console.", output)
			writers = append(writers, getConsoleWriter(&cfg.Zerolog, "stdout"))
		}
	}

	return writers
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

func getFileWriter(rotateCfg *Lumberjack, name string) io.Writer {
	if err := os.MkdirAll(rotateCfg.Directory, 0755); err != nil {
		log.Fatalf("Error during creating the directory (%s) for log file: %s", rotateCfg.Directory, err)
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(rotateCfg.Directory, name),
		MaxSize:    rotateCfg.MaxSize,
		MaxAge:     rotateCfg.MaxAge,
		MaxBackups: rotateCfg.MaxBackups,
		Compress:   rotateCfg.Compress,
	}
}

/*
-------- What I want from logger...
1. I can set up logger with .env variables
2. It includes files rotation (for example with lamberjack)
3. I can easy use this logger in any place of my code (without pass it like a parameter)
4. I don't have deep chain of calls (logger.GetLogger().With().Timestamp().Info().Msg() ...)
5. I can replace gorm logger with my implementation of zerolog logger
*/
