package main

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func init() {
	logger := NewLevelHandler(slog.LevelInfo, os.Stderr)
	slog.SetDefault(logger)
}

func UpdateLogger(level slog.Level) {
	logger := NewLevelHandler(level, os.Stderr)
	slog.SetDefault(logger)
}

func NewLevelHandler(level slog.Level, writer io.Writer) *slog.Logger {
	return slog.New(tint.NewHandler(writer, &tint.Options{
		Level:      level,
		TimeFormat: time.RFC3339,
	}))
}
