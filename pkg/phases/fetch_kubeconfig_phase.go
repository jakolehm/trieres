package phases

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jakolehm/trieres/pkg/cluster"
)

type FetchKubeConfigPhase struct{}

func (p *FetchKubeConfigPhase) Title() string {
	return "Close SSH Connection"
}

func (p *FetchKubeConfigPhase) Run(config *cluster.Config) error {
	masters := config.MasterHosts()
	master := *masters[0]

	configDataBytes, err := master.ExecWithOutput("sudo cat /etc/rancher/k3s/k3s.yaml")
	if err != nil {
		return err
	}
	configData := string(configDataBytes)
	configData = strings.Replace(configData, "https://127.0.0.1:6443", fmt.Sprintf("https://%s:6443", master.Address), 1)

	ioutil.WriteFile("./kubeconfig", []byte(configData), 0700)

	return nil
}
