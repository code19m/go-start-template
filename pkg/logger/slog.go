package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	slogzerolog "github.com/samber/slog-zerolog/v2"
)

const (
	DebugLvl = "debug"
	InfoLvl  = "info"
	WarnLvl  = "warn"
	ErrorLvl = "error"

	TextFormat = "text"
	JsonFormat = "json"
)

var (
	ErrIncorrectLogLevel  = errors.New("incorrect log level")
	ErrIncorrectLogFormat = errors.New("incorrect log format")

	slogLevelMapper = map[string]slog.Level{
		DebugLvl: slog.LevelDebug,
		InfoLvl:  slog.LevelInfo,
		WarnLvl:  slog.LevelWarn,
		ErrorLvl: slog.LevelError,
	}

	zerologLevelMapper = map[string]zerolog.Level{
		DebugLvl: zerolog.DebugLevel,
		InfoLvl:  zerolog.InfoLevel,
		WarnLvl:  zerolog.WarnLevel,
		ErrorLvl: zerolog.ErrorLevel,
	}
)

func NewSlogLogger(level, format string) (*slog.Logger, error) {
	zerologLevel, ok := zerologLevelMapper[level]
	if !ok {
		return nil, ErrIncorrectLogLevel
	}
	zerolog.SetGlobalLevel(zerologLevel)

	zerologLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if format == TextFormat {
		zerologLogger = zerologLogger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	slogLevel, ok := slogLevelMapper[level]
	if !ok {
		return nil, errors.WithStack(ErrIncorrectLogLevel)
	}

	logger := slog.New(slogzerolog.Option{Level: slogLevel, Logger: &zerologLogger}.NewZerologHandler())

	return logger, nil
}
