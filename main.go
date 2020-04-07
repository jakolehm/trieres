package main

import (
	"log"
	"os"

	"github.com/jakolehm/trieres/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Version gets overridden at build time using -X main.Version=$VERSION
var (
	Version = "dev"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}

func main() {
	app := &cli.App{
		Name:    "trieres",
		Version: Version,
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
