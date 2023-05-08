package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

const (
	client = "client"
	server = "server"
	envKey = "LT_TOKEN"
)

type (
	action = func(*cli.Context) error
)

func newApp(local, server action) *cli.App {
	app := cli.NewApp()
	app.Name = "lt"
	app.Usage = "export local port to remote"

	app.Compiled = time.Now()

	app.Authors = []*cli.Author{
		{
			Name:  "wei.xuan",
			Email: "adamweixuan@gmail.com",
		},
	}
	app.Commands = []*cli.Command{
		newCliCmd(local),
		newServCmd(server),
	}

	return app
}

func newCliCmd(action func(*cli.Context) error) *cli.Command {
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
				Name:    "secret,s",
				Aliases: []string{"s"},
				Usage:   "secret",
				EnvVars: []string{envKey},
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
		Action: action,
	}
}

func newServCmd(action func(*cli.Context) error) *cli.Command {
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
				EnvVars: []string{envKey},
				Usage:   "secret",
			},
		},
		Action: action,
	}
}