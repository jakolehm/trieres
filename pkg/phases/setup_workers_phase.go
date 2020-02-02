package phases

import (
	"fmt"
	retry "github.com/avast/retry-go"
	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type SetupWorkersPhase struct{}

var workerSetupCmd = "curl -sfL https://get.k3s.io | %s sh -s - agent --server https://%s:6443 --token %s %s"

func (p *SetupWorkersPhase) Title() string {
	return "Setup k3s workers"
}

func (p *SetupWorkersPhase) Run(config *cluster.Config) error {
	master := config.MasterHosts()[0]
	wg := sync.WaitGroup{}
	for _, host := range config.WorkerHosts() {
		wg.Add(1)
		go p.setupWorker(&wg, host, master.Address, config)
	}
	wg.Wait()

	return nil
}

func (p *SetupWorkersPhase) setupWorker(wg *sync.WaitGroup, host *hosts.Host, master string, config *cluster.Config) error {
	defer wg.Done()

	err := retry.Do(
		func() error {
			logrus.Infof("%s: setting up k3s worker", host.Address)
			setupCmd := fmt.Sprintf(workerSetupCmd, config.SetupEnvs(), master, config.Token, strings.Join(host.ExtraArgs, " "))
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

	logrus.Printf("%s: k3s worker setup succeeded", host.Address)
	return nil
}
