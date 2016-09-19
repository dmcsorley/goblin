// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerRunStepRequiresImage(t *testing.T) {
	_, err := newRunStep(
		0,
		&config.StepRecord{Type: DockerRunStepType},
	)

	if err == nil {
		t.Error("docker-run step should have failed with no image")
	}
}
