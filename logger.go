package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.StampMilli,
	})
}

func UpdateLogger(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}
