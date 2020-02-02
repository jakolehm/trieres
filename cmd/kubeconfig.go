package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func KubeconfigCommand() *cli.Command {
	kubeconfigFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Specify an alternate cluster YAML file",
			Value:   "cluster.yml",
			EnvVars: []string{"TRIERES_CONFIG"},
		},
	}
	return &cli.Command{
		Name:   "kubeconfig",
		Usage:  "Fetch admin kubeconfig",
		Action: showKubeconfig,
		Flags:  kubeconfigFlags,
	}
}

func showKubeconfig(ctx *cli.Context) error {
	cluster := cluster.Config{}
	configBuffer, configFile, err := resolveClusterFile(ctx)
	if err != nil {
		return err
	}
	logrus.Debugf("Loading config file: %s", configFile)
	cluster.FromYaml(configBuffer)

	master := cluster.MasterHosts()[0]
	if master == nil {
		return errors.New("No masters found")
	}
	err = master.Connect()
	if err != nil {
		return err
	}
	outputBytes, err := master.ExecWithOutput("sudo cat /etc/rancher/k3s/k3s.yaml", nil)
	if err != nil {
		return errors.New("Cannot find kubeconfig")
	}
	output := string(outputBytes)
	output = strings.Replace(output, "https://127.0.0.1:6443", fmt.Sprintf("https://%s:6443", master.Address), 1)
	fmt.Print(output)
	master.Disconnect()

	return nil
}
