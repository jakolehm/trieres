package phases

import (
	"fmt"
	retry "github.com/avast/retry-go"
	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
	"strings"
)

type SetupMastersPhase struct{}

var masterSetupCmd = "sh -c 'curl -sfL https://get.k3s.io | %s sh -s - server --agent-token %s %s'"

func (p *SetupMastersPhase) Title() string {
	return "Setup k3s masters"
}

func (p *SetupMastersPhase) Run(config *cluster.Config) error {
	for _, host := range config.MasterHosts() {
		p.setupMaster(host, config)
	}

	return nil
}

func (p *SetupMastersPhase) setupMaster(host *hosts.Host, config *cluster.Config) error {
	err := retry.Do(
		func() error {
			logrus.Infof("%s: setting up k3s master", host.Address)
			setupCmd := fmt.Sprintf(masterSetupCmd, config.SetupEnvs(), config.Token, strings.Join(host.ExtraArgs, " "))
			err := host.Exec(setupCmd)
			if err != nil {
				logrus.Errorf("%s: failed -> %s", host.Address, err.Error())
			}
			return err
		},
	)
	if err != nil {
		logrus.Errorf("%s: failed to setup k3s", host.Address)
		return err
	}

	logrus.Printf("%s: k3s setup succeeded", host.Address)
	return nil
}
