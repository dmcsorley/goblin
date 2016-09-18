// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerBuildStepRequiresImage(t *testing.T) {
	_, err := newBuildStep(
		0,
		&config.StepRecord{Type:DockerBuildStepType},
	)

	if err == nil {
		t.Error("docker-build step should have failed with no image")
	}
}
