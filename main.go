package main

import (
	"log"
	"os"

	"github.com/jakolehm/trieres/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// VERSION gets overridden at build time using -X main.VERSION=$VERSION
var VERSION = "dev"

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	app := &cli.App{
		Name:    "trieres",
		Version: VERSION,
		Usage:   "k3s cluster lifecycle management tool",
		Commands: []*cli.Command{
			cmd.UpCommand(),
			cmd.KubeconfigCommand(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug,d",
				Usage: "Debug logging",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
