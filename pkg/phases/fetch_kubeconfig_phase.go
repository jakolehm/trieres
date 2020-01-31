package phases

import (
	"io/ioutil"

	"github.com/jakolehm/trieres/pkg/cluster"
)

type FetchKubeConfigPhase struct{}

func (p *FetchKubeConfigPhase) Title() string {
	return "Close SSH Connection"
}

func (p *FetchKubeConfigPhase) Run(config *cluster.Config) error {
	masters := config.MasterHosts()
	master := *masters[0]

	configData, err := master.ExecWithOutput("sudo cat /etc/rancher/k3s/k3s.yaml")
	if err != nil {
		return err
	}

	ioutil.WriteFile("./kubeconfig", configData, 0700)

	return nil
}
