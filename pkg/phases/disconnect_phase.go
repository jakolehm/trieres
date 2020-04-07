package phases

import (
	"sync"

	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
)

type DisconnectPhase struct{}

func (p *DisconnectPhase) Title() string {
	return "Close SSH Connection"
}

func (p *DisconnectPhase) Run(config *cluster.Config) error {
	var wg sync.WaitGroup
	for _, host := range config.Hosts {
		wg.Add(1)
		go p.disconnectHost(host, &wg)
	}
	wg.Wait()

	return nil
}

func (p *DisconnectPhase) disconnectHost(host *hosts.Host, wg *sync.WaitGroup) error {
	defer wg.Done()
	host.Connect()
	logrus.Printf("%s: SSH connection closed", host.FullAddress())
	return nil
}
