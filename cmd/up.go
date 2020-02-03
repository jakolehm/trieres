package cmd

import (
	"fmt"

	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/phases"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func UpCommand() *cli.Command {
	upFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Specify an alternate cluster YAML file",
			Value:   "cluster.yml",
			EnvVars: []string{"TRIERES_CONFIG"},
		},
	}
	return &cli.Command{
		Name:   "up",
		Usage:  "Install or upgrade the cluster",
		Action: clusterUp,
		Flags:  upFlags,
	}
}

func clusterUp(ctx *cli.Context) error {
	fmt.Printf("~~ Trieres (version %s) ~~\n\n", ctx.App.Version)

	cluster := cluster.Config{}
	configBuffer, configFile, err := resolveClusterFile(ctx)
	if err != nil {
		return err
	}
	logrus.Infof("Loading config file: %s", configFile)
	cluster.FromYaml(configBuffer)
	if cluster.Token == "" {
		random, err := GenerateRandomString(16)
		if err != nil {
			return err
		}
		cluster.Token = random
	}

	phaseManager := phases.NewManager(&cluster)

	phaseManager.AddPhase(&phases.ValidateConfigurationPhase{})
	phaseManager.AddPhase(&phases.ConnectPhase{})
	phaseManager.AddPhase(&phases.SetupMastersPhase{})
	phaseManager.AddPhase(&phases.SetupWorkersPhase{})
	phaseManager.AddPhase(&phases.DisconnectPhase{})

	phaseErr := phaseManager.Run()
	if phaseErr != nil {
		return phaseErr
	}

	return nil
}
