package cluster

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/jakolehm/trieres/pkg/hosts"
	"gopkg.in/yaml.v2"
)

// Config describes cluster.yml
type Config struct {
	Hosts     hosts.Hosts `validate:"required,dive,required,gt=0"`
	Token     string      `validate:"omitempty,gt=12"`
	Manifests []string
	Version   string
}

// FromYaml parses config from YAML
func (c *Config) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func (c *Config) Validate() error {
	validator := validator.New()
	return validator.Struct(c)
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

func (c *Config) SetupEnvs() string {
	var envs = []string{}

	if c.Version != "" {
		envs = append(envs, fmt.Sprintf("INSTALL_K3S_VERSION=%s", c.Version))
	}

	return strings.Join(envs, " ")
}
