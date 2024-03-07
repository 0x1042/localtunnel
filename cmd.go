package main

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	uuid.EnableRandPool()

	app := cli.Command{}
	app.Name = "localtunnel"
	app.Usage = "export local port to public server"
	app.Version = Version + ". build at " + Date

	app.Flags = []cli.Flag{verbos()}

	app.Commands = []*cli.Command{
		newCliCmd(),
		newServCmd(),
	}

	app.Action = func(ctx context.Context, cmd *cli.Command) error {
		cli.DefaultAppComplete(ctx, cmd)
		cli.ShowVersion(cmd)
		return nil
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("start error.")
	}
}
