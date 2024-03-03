package main

import (
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	uuid.EnableRandPool()

	local := func(ctx *cli.Context) error {
		local := ctx.Int("local")
		remote := ctx.Int("remote")
		tunnel := ctx.String("tunnel")
		secret := ctx.String("secret")
		return StartClient(local, remote, secret, tunnel)
	}

	server := func(ctx *cli.Context) error {
		port := ctx.Int("port")
		secret := ctx.String("secret")
		return StartServer(port, secret)
	}

	if err := newApp(local, server).Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("start error.")
	}
}
