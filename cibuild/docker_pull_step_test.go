// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerPullStepRequiresImage(t *testing.T) {
	expectStepConstructorFailure(
		newPullStep,
		&config.StepRecord{Type: string(DockerPull)},
		t,
		"docker-pull step should have failed with no image",
	)
}

func TestDockerPullStepFailsForBadImageValue(t *testing.T) {
	expectStepConstructorFailure(
		newPullStep,
		&config.StepRecord{
			Type:          string(DockerPull),
			Image:         "${badexample}",
			DecodedFields: []string{"image"},
		},
		t,
		"docker-pull step should have failed for bad image value",
	)
}
