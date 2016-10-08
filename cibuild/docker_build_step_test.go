// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerBuildStepRequiresImage(t *testing.T) {
	expectStepConstructorFailure(
		newBuildStep,
		&config.StepRecord{Type: string(DockerBuild)},
		t,
		"docker-build step should have failed with no image",
	)
}

func TestDockerBuildStepFailsForBadImage(t *testing.T) {
	expectStepConstructorFailure(
		newBuildStep,
		&config.StepRecord{
			Type:          string(DockerBuild),
			Image:         "${badexample}",
			DecodedFields: []string{"image"},
		},
		t,
		"docker-build step should have failed with no image",
	)
}
