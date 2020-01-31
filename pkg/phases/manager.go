package phases

import (
	"github.com/jakolehm/trieres/pkg/cluster"
	"github.com/sirupsen/logrus"
)

type PhaseManager struct {
	phases []Phase
	config *cluster.Config
}

type PhaseContext struct {
	token string
}

const (
	contextConfig int = 1
)

func NewManager(config *cluster.Config) *PhaseManager {
	manager := PhaseManager{
		config: config,
	}
	return &manager
}

// AddPhase adds a Phase to PhaseManager
func (m *PhaseManager) AddPhase(phase Phase) {
	m.phases = append(m.phases, phase)
}

// Run executes all the added Phases in order
func (m *PhaseManager) Run() error {
	for _, phase := range m.phases {
		logrus.Infof("==> Running phase: %s", phase.Title())
		err := phase.Run(m.config)
		if err != nil {
			return err
		}
	}

	return nil
}
