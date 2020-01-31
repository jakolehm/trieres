package cluster

import (
	"github.com/jakolehm/trieres/pkg/hosts"
	"gopkg.in/yaml.v2"
)

// Config describes cluster.yml
type Config struct {
	Hosts hosts.Hosts
	Token string
}

// FromYaml parses config from YAML
func (c *Config) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func (c *Config) MasterHosts() hosts.Hosts {
	masters := hosts.Hosts{}
	for _, host := range c.Hosts {
		if host.Role == "master" {
			masters = append(masters, host)
		}
	}

	return masters
}

func (c *Config) WorkerHosts() hosts.Hosts {
	workers := hosts.Hosts{}
	for _, host := range c.Hosts {
		if host.Role == "worker" {
			workers = append(workers, host)
		}
	}

	return workers
}
