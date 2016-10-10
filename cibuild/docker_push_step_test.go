// import github.com/dmcsorley/goblin/cibuild
package cibuild

import (
	"github.com/dmcsorley/goblin/config"
	"testing"
)

func TestDockerPushStepParams(t *testing.T) {
	dps := &config.StepRecord{Type: string(DockerPush)}

	expectStepConstructorFailure(
		newPushStep,
		dps,
		t,
		"docker-push step should have failed with no image",
	)

	dps.Image = "${badexample}"
	dps.DecodedFields = []string{"image"}

	expectStepConstructorFailure(
		newPushStep,
		dps,
		t,
		"docker-push step should have failed for bad image value",
	)
}
