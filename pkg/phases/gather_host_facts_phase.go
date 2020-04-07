package phases

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
)

type GatherHostFactsPhase struct{}

func (p *GatherHostFactsPhase) Title() string {
	return "Gather Host Facts"
}

func (p *GatherHostFactsPhase) Run(config *cluster.Config) error {
	var wg sync.WaitGroup
	for _, host := range config.Hosts {
		wg.Add(1)
		logrus.Infof("%s: gathering facts", host.FullAddress())
		go investigateHost(host, &wg)
	}
	wg.Wait()

	return nil
}

func investigateHost(host *hosts.Host, wg *sync.WaitGroup) error {
	defer wg.Done()

	orgName := host.FullAddress()
	host.Metadata = &hosts.HostMetadata{
		Hostname:        resolveHostname(host),
		InternalAddress: resolveInternalIP(host),
	}
	if host.Metadata.Hostname != "" {
		logrus.Infof("%s: is now known as %s", orgName, host.FullAddress())
	}
	return nil
}

func resolveHostname(host *hosts.Host) string {
	hostname, _ := host.ExecWithOutput("hostname -s")

	return hostname
}

func resolveInternalIP(host *hosts.Host) string {
	//ip -o addr show dev #{interface} scope global
	output, _ := host.ExecWithOutput(fmt.Sprintf("ip -o addr show dev %s scope global", host.PrivateInterface))
	//logrus.Debugln(output)
	lines := strings.Split(output, "\r\n")
	for _, line := range lines {
		items := strings.Fields(line)
		addrItems := strings.Split(items[3], "/")
		if addrItems[0] != host.Address {
			return addrItems[0]
		}
	}
	return host.Address
}
