package phases

import (
	"fmt"
	"strings"

	retry "github.com/avast/retry-go"
	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
)

type SetupMastersPhase struct{}

var masterSetupCmd = "sh -c 'curl -sfL https://get.k3s.io | %s sh -s - server --agent-token \"%s\" %s'"

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
	if strings.HasPrefix(host.Address, "127.") && host.Metadata.Hostname != "" {
		host.ExtraArgs = append(host.ExtraArgs, fmt.Sprintf("--node-ip=%s", host.Metadata.InternalAddress))
	}
	err := retry.Do(
		func() error {
			logrus.Infof("%s: setting up k3s master", host.FullAddress())
			setupCmd := fmt.Sprintf(masterSetupCmd, config.SetupEnvs(), config.Token, strings.Join(host.ExtraArgs, " "))
			err := host.Exec(setupCmd)
			if err != nil {
				logrus.Errorf("%s: failed -> %s", host.FullAddress(), err.Error())
			}
			return err
		},
	)
	if err != nil {
		logrus.Errorf("%s: failed to setup k3s", host.FullAddress())
		return err
	}

	logrus.Printf("%s: k3s setup succeeded", host.FullAddress())
	return nil
}
