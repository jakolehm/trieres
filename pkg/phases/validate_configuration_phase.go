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
		if host.Role != "master" && host.Role != "worker" {
			messages = append(messages, fmt.Sprintf("%s: Invalid role: %s", host.Address, host.Role))
		}

		if host.SSHPort < 1 || host.SSHPort > 65535 {
			messages = append(messages, fmt.Sprintf("%s: Invalid SSH Port: %d", host.Address, host.SSHPort))
		}

		for _, host2 := range config.Hosts {
			if host2.Address == host.Address && host2.SSHPort == host.SSHPort {
				messages = p.AppendUnlessContains(messages, fmt.Sprintf("Duplicate address:ssh_port %s:%d", host.Address, host.SSHPort))
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

func (p *ValidateConfigurationPhase) AppendUnlessContains(list []string, item string) []string {
	messages := []string{}
	if p.ContainsString(list, item) {
		return list
	}
	messages = append(list, item)
	return messages
}

func (p *ValidateConfigurationPhase) ContainsString(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}
