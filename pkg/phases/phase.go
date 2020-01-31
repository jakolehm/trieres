package phases

import (
	"github.com/jakolehm/trieres/pkg/cluster"
)

// Phase interface
type Phase interface {
	Run(config *cluster.Config) error
	Title() string
}
