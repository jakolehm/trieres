package phases

import (
	"sync"

	retry "github.com/avast/retry-go"
	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
)

type ConnectPhase struct{}

func (p *ConnectPhase) Title() string {
	return "Open SSH Connection"
}

func (p *ConnectPhase) Run(config *cluster.Config) error {
	var wg sync.WaitGroup
	for _, host := range config.Hosts {
		wg.Add(1)
		go p.connectHost(host, &wg)
	}
	wg.Wait()

	return nil
}

func (p *ConnectPhase) connectHost(host *hosts.Host, wg *sync.WaitGroup) error {
	defer wg.Done()
	err := retry.Do(
		func() error {
			logrus.Infof("%s: opening SSH connection", host.Address)
			err := host.Connect()
			if err != nil {
				logrus.Errorf("%s: failed to connect -> %s", host.Address, err.Error())
			}
			return err
		},
	)
	if err != nil {
		logrus.Errorf("%s: failed to open connection", host.Address)
		return err
	}

	logrus.Printf("%s: SSH connection opened", host.Address)
	return nil
}
