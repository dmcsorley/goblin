// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerPullStepRequiresImage(t *testing.T) {
	_, err := newPullStep(
		0,
		&config.StepRecord{Type: DockerPullStepType},
	)

	if err == nil {
		t.Error("docker-pull step should have failed with no image")
	}
}
