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
		NoColor:    false,
		TimeFormat: time.StampMilli,
	})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func enableTrace() error {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	return nil
}
