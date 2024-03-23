package main

import (
	"context"

	"github.com/urfave/cli/v3"
)

const (
	client = "client"
	server = "server"
	envKey = "LT_TOKEN"
)

func verbos() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "verbos",
		Value: false,
		Usage: "verbos log",
		Action: func(_ context.Context, _ *cli.Command, _ bool) error {
			return enableTrace()
		},
	}
}

func newCliCmd() *cli.Command {
	return &cli.Command{
		Name:  client,
		Usage: "start tunnel client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tunnel",
				Aliases:  []string{"t"},
				Usage:    "tunnel server addr",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "secret",
				Sources: cli.EnvVars(envKey),
			},
			&cli.IntFlag{
				Name:     "local",
				Required: true,
				Aliases:  []string{"l"},
				Usage:    "local port",
			},
			&cli.IntFlag{
				Name:        "remote",
				Aliases:     []string{"r"},
				Value:       0,
				Usage:       "remote port",
				DefaultText: "random",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			local := cmd.Int("local")
			remote := cmd.Int("remote")
			secret := cmd.String("secret")
			tunnel := cmd.String("tunnel")
			return StartClient(int(local), int(remote), secret, tunnel)
		},
	}
}

func newServCmd() *cli.Command {
	return &cli.Command{
		Name:  server,
		Usage: "start tunnel server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   7853,
				Usage:   "tunnel listen port",
			},
			&cli.StringFlag{
				Name:    "secret",
				Aliases: []string{"s"},
				Sources: cli.EnvVars(envKey),
				Usage:   "secret",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			port := cmd.Int("port")
			secret := cmd.String("secret")
			return StartServer(int(port), secret)
		},
	}
}
