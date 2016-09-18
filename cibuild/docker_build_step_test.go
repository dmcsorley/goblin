package cibuild

import (
	"testing"
)

func TestDockerBuildStepRequiresImage(t *testing.T) {
	_, err := newBuildStep(
		0,
		map[string]interface{}{
			TypeKey: DockerBuildStepType,
		},
	)

	if err == nil {
		t.Error("docker-build step should have failed with no image")
	}
}
