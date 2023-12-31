package prepare

import (
	"log/slog"
	"os"
)

const defaultLogLevel = slog.LevelInfo

func SLogger() *slog.Logger {

	opts := &slog.HandlerOptions{
		Level: defaultLogLevel,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	return logger
}
