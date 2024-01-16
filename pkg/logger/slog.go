package logger

import (
	"log/slog"

	"github.com/pkg/errors"
	slogzap "github.com/samber/slog-zap"
	"go.uber.org/zap"
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
)

func NewSlogLogger(level, format string) (*slog.Logger, error) {

	zapLogger, err := zap.NewDevelopment()

	if err != nil {
		return nil, errors.Wrap(err, "failed to init zap logger")
	}

	slogLevel, ok := slogLevelMapper[level]
	if !ok {
		return nil, errors.WithStack(ErrIncorrectLogLevel)
	}

	logger := slog.New(slogzap.Option{Level: slogLevel, Logger: zapLogger}.NewZapHandler())

	// TODO: Fix bug with debug level

	return logger, nil
}
