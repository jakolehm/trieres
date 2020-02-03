package phases

import (
	"errors"
	"fmt"
	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/sirupsen/logrus"
)

type ValidateConfigurationPhase struct{}

func (p *ValidateConfigurationPhase) Title() string {
	return "Validate Configuration"
}

func (p *ValidateConfigurationPhase) Run(config *cluster.Config) error {
	messages := []string{}

	for _, host := range config.Hosts {
		logrus.Infof("%s: Validating", host.Address)
		for _, host2 := range config.Hosts {
			if host2.Address == host.Address && host2.SSHPort == host.SSHPort {
				message := fmt.Sprintf("Duplicate address:ssh_port %s:%d", host.Address, host.SSHPort)
				if !p.ContainsString(messages, message) {
					messages = append(messages, message)
				}
			}
		}
	}

	if len(messages) > 0 {
		for _, message := range messages {
			logrus.Error(message)
		}
		return errors.New("Invalid configuration")
	}

	return nil
}

func (p *ValidateConfigurationPhase) ContainsString(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}
