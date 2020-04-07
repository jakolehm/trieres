package phases

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/jakolehm/trieres/pkg/hosts"
	"github.com/sirupsen/logrus"
)

var hostPath = "/var/lib/rancher/k3s/server/manifests/trieres"

type CopyManifestsPhase struct{}

// Title returns phase title
func (p *CopyManifestsPhase) Title() string {
	return "Copy manifests"
}

// Run executes phase
func (p *CopyManifestsPhase) Run(config *cluster.Config) error {
	masterHost := config.MasterHosts()[0]

	err := masterHost.Exec(fmt.Sprintf("sudo mkdir -p %s", hostPath))
	if err != nil {
		return err
	}
	masterHost.Exec(fmt.Sprintf("sudo rm -f %s/*.yaml", hostPath))

	for _, manifestGlob := range config.Manifests {
		manifests, err := filepath.Glob(manifestGlob)
		if err != nil {
			return err
		}
		for _, manifest := range manifests {
			file, err := os.Open(manifest)
			if err != nil {
				return err
			}
			logrus.Infof("%s: copying manifest %s", masterHost.FullAddress(), file.Name())
			err = handleManifest(masterHost, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func handleManifest(masterHost *hosts.Host, file *os.File) error {
	hasher := md5.New()
	hasher.Write([]byte(file.Name()))
	targetFile := path.Join(hostPath, fmt.Sprintf("%s.yaml", hex.EncodeToString(hasher.Sum(nil))))
	copyErr := masterHost.CopyFile(*file, targetFile, "0600")
	if copyErr != nil {
		return copyErr
	}

	return nil
}
