package main

import (
	"os"
	"time"

	"github.com/lmittmann/tint"
	"golang.org/x/exp/slog"
)

func init() {
	level := slog.LevelInfo
	if os.Getenv("verbos") == "1" {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
		TimeFormat: time.RFC3339,
	})))
}
